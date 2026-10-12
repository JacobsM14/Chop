package store

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"testing"
)

func TestOpen(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestCreateChopDuplicate(t *testing.T) {
	s := openTest(t)
	createChop(t, s, "deploy")

	if err := s.CreateChop("deploy", "red", false); !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
}

func TestItemPositions(t *testing.T) {
	s := openTest(t)
	createChop(t, s, "deploy", "a", "b", "c")

	// Removing an item leaves a gap, and new items go after the last one.
	if err := s.RmItem("deploy", 2); err != nil {
		t.Fatal(err)
	}
	if err := s.AddItem("deploy", ListItem{Text: "d"}); err != nil {
		t.Fatal(err)
	}

	chop, err := s.Get("deploy")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, it := range chop.Items {
		got = append(got, fmt.Sprintf("%d %s", it.Position, it.Text))
	}
	if want := []string{"1 a", "3 c", "4 d"}; !slices.Equal(got, want) {
		t.Errorf("items = %q, want %q", got, want)
	}
}

func TestDeleteChopRemovesItems(t *testing.T) {
	s := openTest(t)
	createChop(t, s, "deploy", "a", "b")
	if err := s.DeleteChop("deploy"); err != nil {
		t.Fatal(err)
	}

	// SQLite reuses the freed id, so if the items survived the delete the
	// new chop would show them.
	createChop(t, s, "deploy")
	chop, err := s.Get("deploy")
	if err != nil {
		t.Fatal(err)
	}
	if len(chop.Items) != 0 {
		t.Errorf("recreated chop has %d items, want 0", len(chop.Items))
	}
}

func TestNotFound(t *testing.T) {
	s := openTest(t)
	createChop(t, s, "deploy", "a")

	tests := []struct {
		name string
		fn   func() error
	}{
		{"Get", func() error { _, err := s.Get("nope"); return err }},
		{"AddItem", func() error { return s.AddItem("nope", ListItem{Text: "a"}) }},
		{"DeleteChop", func() error { return s.DeleteChop("nope") }},
		{"RmItem missing chop", func() error { return s.RmItem("nope", 1) }},
		{"RmItem missing position", func() error { return s.RmItem("deploy", 2) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); !errors.Is(err, ErrNotFound) {
				t.Errorf("got %v, want ErrNotFound", err)
			}
		})
	}
}

func TestListInvalidSort(t *testing.T) {
	s := openTest(t)
	if _, err := s.List("size"); err == nil {
		t.Fatal("got no error for an invalid sort")
	}
}

func openTest(t *testing.T) Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// createChop creates a chop with one item per text.
func createChop(t *testing.T, s Store, name string, texts ...string) {
	t.Helper()
	if err := s.CreateChop(name, "green", true); err != nil {
		t.Fatal(err)
	}
	for _, text := range texts {
		if err := s.AddItem(name, ListItem{Text: text}); err != nil {
			t.Fatal(err)
		}
	}
}
