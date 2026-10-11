package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"chop/internal/render"
	"chop/internal/store"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmItemCmd)
}

var rmItemCmd = &cobra.Command{
	Use:   "rm-item <name> [n]",
	Short: "Remove an item from a Chop",
	Long:  "Remove an item from a Chop. <n> is the item number as shown by `chop get`; if omitted, pick the item from a list.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		chop, err := db.Get(name)
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("chop %q not found", name)
		}
		if err != nil {
			return err
		}
		if len(chop.Items) == 0 {
			return fmt.Errorf("chop %q has no items", name)
		}

		var idx int
		if len(args) == 2 {
			n, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid item number %q", args[1])
			}
			if n < 1 || n > len(chop.Items) {
				return fmt.Errorf("chop %q has no item %d (it has %d)", name, n, len(chop.Items))
			}
			idx = n - 1
		} else {
			err := selectItem(chop, &idx)
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("running form: %w", err)
			}
		}

		// The number shown by `chop get` is the item's index, not its stored
		// position: positions keep gaps after items are removed.
		return db.RmItem(name, chop.Items[idx].Position)
	},
}

func selectItem(chop store.Chop, idx *int) error {
	opts := make([]huh.Option[int], len(chop.Items))
	for i, it := range chop.Items {
		opts[i] = huh.NewOption(render.Item(i, chop.Ordered, it), i)
	}

	return huh.NewSelect[int]().
		Title("Item to remove").
		Options(opts...).
		Value(idx).
		Run()
}
