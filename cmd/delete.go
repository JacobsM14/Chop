package cmd

import (
	"errors"
	"fmt"

	"github.com/JacobsM14/chop/internal/store"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var deleteYes bool

func init() {
	deleteCmd.Flags().BoolVarP(&deleteYes, "yes", "y", false, "skip confirmation")
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:               "delete <name>",
	Aliases:           []string{"rm"},
	Short:             "Delete a Chop and all its items",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeChopNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if !deleteYes {
			var confirm bool
			err := huh.NewConfirm().
				Title(fmt.Sprintf("Delete chop %q and all its items?", name)).
				Value(&confirm).
				Run()
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("running form: %w", err)
			}
			if !confirm {
				return nil
			}
		}

		err := db.DeleteChop(name)
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("chop %q not found", name)
		}
		return err
	},
}
