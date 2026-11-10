package spaces

import (
	"context"
	"testing"
)

func TestDeletionStateSurvivesARestart(t *testing.T) {
	db := openTestDatabase(t)
	state := deletionState{db: db}
	if saved, err := state.Load(context.Background()); err != nil || saved != nil {
		t.Fatalf("empty Load = %q, %v; want nil", saved, err)
	}
	for _, blob := range []string{"one", "two"} {
		if err := state.Save(context.Background(), []byte(blob)); err != nil {
			t.Fatalf("save: %v", err)
		}
	}
	if saved, err := state.Load(context.Background()); err != nil || string(saved) != "two" {
		t.Fatalf("Load = %q, %v; want two", saved, err)
	}
}
