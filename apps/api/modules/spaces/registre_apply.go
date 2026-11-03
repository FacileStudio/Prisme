package spaces

import (
	"context"
	stderrors "errors"
	"fmt"
	"slices"

	"github.com/FacileStudio/Prisme/apps/api/internal/matrix"
	"github.com/FacileStudio/Prisme/apps/api/schemas"
	"github.com/FacileStudio/porte/directory"

	"gorm.io/gorm"
)

// Apply makes Prisme's linked spaces match the snapshot: every Registre space
// gets a Prisme space with #general, and members join or leave the space and
// its channels, their agents with them. A space that fails is reported and the
// others still apply. It is idempotent.
func (r *Registre) Apply(ctx context.Context, snapshot directory.Snapshot) error {
	var failures []error
	for _, remote := range snapshot.Spaces {
		if err := r.applySpace(ctx, remote); err != nil {
			failures = append(failures, fmt.Errorf("directory: space %s: %w", remote.ID, err))
		}
	}
	return stderrors.Join(failures...)
}

func (r *Registre) applySpace(ctx context.Context, remote directory.Space) error {
	people, err := r.resolveAll(ctx, remote.Members)
	if err != nil {
		return err
	}
	row, err := r.link(ctx, remote, people)
	if err != nil || row.ID == 0 {
		return err
	}
	return r.reconcile(ctx, row, people)
}

// link returns the Prisme space for a Registre space: the one linked to it,
// else the oldest unlinked one with the same name, case-insensitively, else a
// new one created by the snapshot's owner. A space with no resolvable member
// and no local copy is skipped (an empty row).
func (r *Registre) link(ctx context.Context, remote directory.Space, people []person) (schemas.Space, error) {
	db := r.spaces.db.WithContext(ctx)
	var row schemas.Space
	err := db.Where("registre_space_id = ?", remote.ID).First(&row).Error
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		err = db.Where("registre_space_id IS NULL AND lower(name) = lower(?)", remote.Name).Order("created_at, id").First(&row).Error
	}
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		return r.create(ctx, remote, people)
	}
	if err != nil {
		return row, err
	}
	return row, r.relink(ctx, &row, remote)
}

// relink stores the Registre id and name on the row and renames the Matrix
// space when Registre's name differs.
func (r *Registre) relink(ctx context.Context, row *schemas.Space, remote directory.Space) error {
	if row.Name == remote.Name && row.RegistreSpaceID != nil {
		return nil
	}
	err := r.spaces.db.WithContext(ctx).Model(&schemas.Space{}).Where("id = ?", row.ID).
		Updates(map[string]any{"name": remote.Name, "registre_space_id": remote.ID}).Error
	if err != nil || row.Name == remote.Name {
		return err
	}
	row.Name = remote.Name
	event := matrix.StateEvent{Type: "m.room.name", Content: map[string]string{"name": remote.Name}}
	return r.spaces.matrix.SendState(ctx, r.spaces.bot(), row.RoomID, event)
}

// create provisions the space as its Registre owner, else its first member,
// and brings that creator's agents in.
func (r *Registre) create(ctx context.Context, remote directory.Space, people []person) (schemas.Space, error) {
	if len(people) == 0 {
		return schemas.Space{}, nil
	}
	creator := people[0]
	if at := slices.IndexFunc(people, func(p person) bool { return p.role == directory.RoleOwner }); at >= 0 {
		creator = people[at]
	}
	account, err := r.spaces.accounts.Ensure(ctx, creator.userID)
	if err != nil {
		return schemas.Space{}, err
	}
	remoteID := remote.ID
	space, err := r.spaces.provision(ctx, account, remote.Name, &remoteID)
	if err != nil {
		return schemas.Space{}, err
	}
	r.spaces.follow(ctx, creator.userID, roomsOf(space))
	var row schemas.Space
	return row, r.spaces.db.WithContext(ctx).First(&row, "id = ?", space.ID).Error
}

// reconcile admits the wanted people who are not joined to the space and
// evicts the Prisme accounts joined to it that the snapshot no longer holds.
func (r *Registre) reconcile(ctx context.Context, row schemas.Space, people []person) error {
	rooms, err := r.rooms(ctx, row)
	if err != nil {
		return err
	}
	joined, err := r.spaces.matrix.JoinedMembers(ctx, r.spaces.bot(), row.RoomID)
	if err != nil {
		return fmt.Errorf("list the space members: %w", err)
	}
	current, err := r.accountsIn(ctx, joined)
	if err != nil {
		return err
	}
	wanted := map[int64]bool{}
	var failures []error
	for _, p := range people {
		wanted[p.userID] = true
		if _, in := current[p.userID]; !in {
			failures = append(failures, r.admit(ctx, p.userID, rooms))
		}
	}
	for userID, matrixID := range current {
		if !wanted[userID] {
			failures = append(failures, r.evict(ctx, userID, matrixID, rooms))
		}
	}
	return stderrors.Join(failures...)
}

// admit invites the user into the space as the bot, joins them to the space
// and every channel (restricted to space members), then brings their agents.
func (r *Registre) admit(ctx context.Context, userID int64, rooms []string) error {
	account, err := r.spaces.accounts.Ensure(ctx, userID)
	if err != nil {
		return err
	}
	if err := r.spaces.matrix.Invite(ctx, r.spaces.bot(), rooms[0], account.MatrixUserID); err != nil {
		return fmt.Errorf("invite %s: %w", account.MatrixUserID, err)
	}
	for _, roomID := range rooms {
		if err := r.spaces.matrix.Join(ctx, account.MatrixUserID, roomID); err != nil {
			return fmt.Errorf("join %s to %s: %w", account.MatrixUserID, roomID, err)
		}
	}
	r.spaces.follow(ctx, userID, rooms)
	return nil
}

// evict makes the user leave the channels, then the space, and takes their
// agents out with them.
func (r *Registre) evict(ctx context.Context, userID int64, matrixID string, rooms []string) error {
	joined, err := r.spaces.matrix.JoinedRooms(ctx, matrixID)
	if err != nil {
		return fmt.Errorf("list the rooms of %s: %w", matrixID, err)
	}
	for i := len(rooms) - 1; i >= 0; i-- {
		if !slices.Contains(joined, rooms[i]) {
			continue
		}
		if err := r.spaces.matrix.Leave(ctx, matrixID, rooms[i]); err != nil {
			return fmt.Errorf("%s leaves %s: %w", matrixID, rooms[i], err)
		}
	}
	r.spaces.unfollow(ctx, userID, rooms)
	return nil
}
