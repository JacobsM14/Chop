package cmd

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JacobsM14/chop/internal/store"
)

func TestCompletion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	t.Setenv("CHOP_DB", path)

	s, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"deploy", "dev"} {
		if err := s.CreateChop(name, "green", true); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.AddItem("deploy", store.ListItem{Text: "go test ./..."}); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	// The last line is the directive: 4 = ShellCompDirectiveNoFileComp.
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"run", ""}, "deploy\ndev\n:4\n"},
		{[]string{"get", "deploy", ""}, ":4\n"},
		{[]string{"rm-item", "deploy", ""}, "1\tgo test ./...\n:4\n"},
		{[]string{"rm-item", "nope", ""}, ":4\n"},
		{[]string{"run", "deploy", "--item", ""}, "1\tgo test ./...\n:4\n"},
		{[]string{"new", "--color", ""}, "red\ngreen\nblue\npurple\norange\nwhite\n:4\n"},
		{[]string{"ls", "--sort", ""}, "name\ncolor\nordered\n:4\n"},
	}
	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			var out bytes.Buffer
			rootCmd.SetOut(&out)
			rootCmd.SetErr(io.Discard)
			rootCmd.SetArgs(append([]string{"__complete"}, tt.args...))
			if err := rootCmd.Execute(); err != nil {
				t.Fatal(err)
			}
			if got := out.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
