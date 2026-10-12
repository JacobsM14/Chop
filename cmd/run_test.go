package cmd

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/JacobsM14/chop/internal/store"
)

func TestRun(t *testing.T) {
	ok := []string{"echo 1 >> out.txt", "echo 2 >> out.txt", "echo 3 >> out.txt"}
	failing := []string{"echo 1 >> out.txt", "exit 3", "echo 3 >> out.txt"}

	// Every case runs against a chop called "test" holding items. want is the
	// content of out.txt after the run, which the items append to.
	tests := []struct {
		name    string
		items   []string
		args    []string
		want    string
		wantErr string
	}{
		{"all items in order", ok, []string{"test"}, "1\n2\n3\n", ""},
		{"single item", ok, []string{"test", "--item", "2"}, "2\n", ""},
		{"item too big", ok, []string{"test", "--item", "4"}, "", `chop "test" has no item 4 (it has 3)`},
		{"item zero", ok, []string{"test", "--item", "0"}, "", `chop "test" has no item 0 (it has 3)`},
		{"dry run", ok, []string{"test", "--dry-run"}, "", ""},
		{"unknown chop", ok, []string{"nope"}, "", `chop "nope" not found`},
		{"empty chop", nil, []string{"test"}, "", ""},
		{"stops on failure", failing, []string{"test"}, "1\n", "item 2 failed: exit status 3"},
		{"continue on failure", failing, []string{"test", "--continue"}, "1\n3\n", "1 of 3 commands failed"},
		{"single item keeps its number", failing, []string{"test", "--item", "2"}, "", "item 2 failed: exit status 3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Chdir(dir)
			setupRunChop(t, filepath.Join(dir, "test.db"), tt.items)
			resetRunFlags(t)

			rootCmd.SetArgs(append([]string{"run", "--yes"}, tt.args...))
			err := rootCmd.Execute()

			switch {
			case tt.wantErr == "" && err != nil:
				t.Fatalf("unexpected error: %v", err)
			case tt.wantErr != "" && err == nil:
				t.Fatalf("got no error, want %q", tt.wantErr)
			case tt.wantErr != "" && err.Error() != tt.wantErr:
				t.Fatalf("got error %q, want %q", err, tt.wantErr)
			}

			out, err := os.ReadFile(filepath.Join(dir, "out.txt"))
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				t.Fatal(err)
			}
			if got := string(out); got != tt.want {
				t.Errorf("out.txt = %q, want %q", got, tt.want)
			}
		})
	}
}

func setupRunChop(t *testing.T, path string, items []string) {
	t.Helper()
	t.Setenv("CHOP_DB", path)

	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.CreateChop("test", "green", true); err != nil {
		t.Fatal(err)
	}
	for _, text := range items {
		if err := s.AddItem("test", store.ListItem{Text: text}); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
}

// resetRunFlags undoes the flags set by a previous Execute: cobra keeps both
// the values and Changed between runs, and run relies on Changed("item").
func resetRunFlags(t *testing.T) {
	t.Helper()
	for _, name := range []string{"item", "dry-run", "continue", "yes"} {
		f := runCmd.Flags().Lookup(name)
		if err := f.Value.Set(f.DefValue); err != nil {
			t.Fatal(err)
		}
		f.Changed = false
	}
}
