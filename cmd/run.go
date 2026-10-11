package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"chop/internal/render"
	"chop/internal/runner"
	"chop/internal/store"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	runItem     int
	runDryRun   bool
	runContinue bool
	runYes      bool
)

func init() {
	runCmd.Flags().IntVarP(&runItem, "item", "i", 0, "run only item N (1-based)")
	runCmd.Flags().BoolVar(&runDryRun, "dry-run", false, "print the commands without running them")
	runCmd.Flags().BoolVar(&runContinue, "continue", false, "keep going if a command fails")
	runCmd.Flags().BoolVarP(&runYes, "yes", "y", false, "skip confirmation")
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run <name>",
	Short: "Run the commands of a Chop in order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		chop, err := db.Get(name)
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("chop %q not found", name)
		}
		if err != nil {
			return err
		}

		// first is the index of items[0] in chop.Items, so messages can use
		// the same numbers as `chop get` when --item narrows the list.
		items, first := chop.Items, 0
		if cmd.Flags().Changed("item") {
			if runItem < 1 || runItem > len(items) {
				return fmt.Errorf("chop %q has no item %d (it has %d)", name, runItem, len(items))
			}
			first = runItem - 1
			items = items[first : first+1]
		}
		if len(items) == 0 {
			fmt.Printf("nothing to run, add an item with: chop add %s\n", name)
			return nil
		}

		if !runYes && !runDryRun {
			ok, err := confirmRun(name, items)
			if errors.Is(err, huh.ErrUserAborted) {
				return nil
			}
			if err != nil {
				return fmt.Errorf("running form: %w", err)
			}
			if !ok {
				return nil
			}
		}

		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer stop()

		failed := 0
		for i, it := range items {
			n := first + i + 1
			fmt.Println(render.Command(it))
			if runDryRun {
				continue
			}

			err := runner.Run(ctx, it.Text)
			if ctx.Err() != nil {
				return errors.New("interrupted")
			}
			if err == nil {
				continue
			}
			if !runContinue {
				return fmt.Errorf("item %d failed: %w", n, err)
			}
			fmt.Fprintf(os.Stderr, "item %d failed: %v\n", n, err)
			failed++
		}

		if failed > 0 {
			return fmt.Errorf("%d of %d commands failed", failed, len(items))
		}
		return nil
	},
}

func confirmRun(name string, items []store.ListItem) (bool, error) {
	lines := make([]string, len(items))
	for i, it := range items {
		lines[i] = render.Command(it)
	}

	var confirm bool
	err := huh.NewConfirm().
		Title(fmt.Sprintf("Run these commands from %q?", name)).
		Description(strings.Join(lines, "\n")).
		Value(&confirm).
		Run()
	return confirm, err
}
