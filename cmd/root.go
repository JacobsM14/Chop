package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"chop/internal/store"

	"github.com/spf13/cobra"
)

var db store.Store

var rootCmd = &cobra.Command{
	Use:           "chop",
	Short:         "Save and run lists of commands",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		path, err := dbPath()
		if err != nil {
			return err
		}
		db, err = store.Open(path)
		return err
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return db.Close()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func dbPath() (string, error) {
	if db := os.Getenv("CHOP_DB"); db != "" {
		return db, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "chop", "chop.db"), nil
}
