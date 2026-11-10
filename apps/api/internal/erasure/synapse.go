package erasure

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
)

// eraseMatrix erases every Matrix account of the person: their own and their
// agents'. Rooms are settled first, while the person is still joined and can
// still hand their power on, then each account loses its media and is
// deactivated with erasure. On a retry an account already deactivated is in
// no room any more, and every call repeats harmlessly.
func (e *Eraser) eraseMatrix(ctx context.Context, userID int64) error {
	ids, err := e.matrixIDs(ctx, userID)
	if err != nil || len(ids) == 0 {
		return err
	}
	if err := e.settleRooms(ctx, ids); err != nil {
		return err
	}
	for _, id := range ids {
		if err := e.deps.Matrix.DeleteUserMedia(ctx, id); err != nil {
			return fmt.Errorf("erasure: delete the media of an account: %w", err)
		}
		if err := e.deps.Matrix.EraseUser(ctx, id); err != nil {
			return fmt.Errorf("erasure: deactivate an account: %w", err)
		}
	}
	return nil
}

func (e *Eraser) matrixIDs(ctx context.Context, userID int64) ([]string, error) {
	var ids []string
	err := e.deps.DB.WithContext(ctx).Raw(personMatrixIDs, userID).Scan(&ids).Error
	slices.Sort(ids)
	return ids, err
}

// settleRooms applies the room rule to every room an account of the person
// is joined to: purge a room nobody else is in, unless Prisme lists it as a
// space or channel; keep every shared room, and when no one left in it could
// still administer it, raise the most powerful remaining member first.
func (e *Eraser) settleRooms(ctx context.Context, ids []string) error {
	rooms, err := e.joinedRooms(ctx, ids)
	if err != nil {
		return err
	}
	tracked, err := e.trackedRooms(ctx)
	if err != nil {
		return err
	}
	erased := map[string]bool{}
	for _, id := range ids {
		erased[id] = true
	}
	for _, roomID := range rooms {
		if err := e.settleRoom(ctx, roomID, erased, tracked[roomID]); err != nil {
			return fmt.Errorf("erasure: settle room %s: %w", roomID, err)
		}
	}
	return nil
}

func (e *Eraser) settleRoom(ctx context.Context, roomID string, erased map[string]bool, tracked bool) error {
	members, err := e.deps.Matrix.RoomMembers(ctx, roomID)
	if err != nil {
		return err
	}
	state, err := e.deps.Matrix.RoomState(ctx, roomID)
	if err != nil {
		return err
	}
	verdict := decide(room{members: members, power: readPower(state), tracked: tracked}, erased, e.deps.Bot)
	switch verdict.action {
	case actionPurge:
		return e.deps.Matrix.PurgeRoom(ctx, roomID)
	case actionPromote:
		return e.deps.Matrix.MakeRoomAdmin(ctx, roomID, verdict.promote)
	}
	if verdict.unsettled != "" {
		e.deps.Logger.Warn("erasure: left a room as it was", slog.String("room", roomID), slog.String("reason", verdict.unsettled))
	}
	return nil
}

func (e *Eraser) joinedRooms(ctx context.Context, ids []string) ([]string, error) {
	var rooms []string
	for _, id := range ids {
		joined, err := e.deps.Matrix.UserRooms(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("erasure: list the rooms of an account: %w", err)
		}
		rooms = append(rooms, joined...)
	}
	slices.Sort(rooms)
	return slices.Compact(rooms), nil
}

func (e *Eraser) trackedRooms(ctx context.Context) (map[string]bool, error) {
	var ids []string
	err := e.deps.DB.WithContext(ctx).Raw(`SELECT room_id FROM spaces UNION SELECT room_id FROM channels`).Scan(&ids).Error
	tracked := make(map[string]bool, len(ids))
	for _, id := range ids {
		tracked[id] = true
	}
	return tracked, err
}
