package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/JacobsM14/chop/internal/render"
	"github.com/JacobsM14/chop/internal/store"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var (
	newColor   string
	newOrdered bool
)

// Initialize the new command and its flags
func init() {
	newCmd.Flags().StringVar(&newColor, "color", "", "color principal")
	newCmd.Flags().BoolVar(&newOrdered, "ordered", false, "lista numerada")
	cobra.CheckErr(newCmd.RegisterFlagCompletionFunc("color", cobra.FixedCompletions(render.Names(), cobra.ShellCompDirectiveNoFileComp)))
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new Chop",
	Args:  cobra.ExactArgs(1),
	// The name is new, so there is nothing to complete.
	ValidArgsFunction: cobra.NoFileCompletions,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		var err error
		// If the --color or --ordered flags are not provided, prompt the user for input using the form.
		switch {
		case !cmd.Flags().Changed("color") && !cmd.Flags().Changed("ordered"):
			err = formBoth(&newColor, &newOrdered)
		case !cmd.Flags().Changed("color"):
			err = formColor(&newColor)
		}
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("running form: %w", err)
		}

		if !render.Valid(newColor) {
			return fmt.Errorf("invalid color %q (valid: %s)", newColor, strings.Join(render.Names(), ", "))
		}

		// Here you would typically call the function to create the new Chop in the store.
		// For example:
		err = db.CreateChop(name, newColor, newOrdered)
		if errors.Is(err, store.ErrAlreadyExists) {
			return fmt.Errorf("chop %q already exists", name)
		}
		if err != nil {
			return fmt.Errorf("creating chop: %w", err)
		}
		return nil
	},
}

func colorOptions() []huh.Option[string] {
	opts := make([]huh.Option[string], 0, len(render.Colors))
	for _, c := range render.Colors {
		opts = append(opts, huh.NewOption(render.Paint(c.Name, c.Name), c.Name))
	}
	return opts
}
func formBoth(color *string, ordered *bool) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Color").
				Options(colorOptions()...).
				Value(color),

			huh.NewConfirm().
				Title("Ordered list?").
				Value(ordered),
		),
	)

	return form.Run()
}

func formColor(color *string) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Color").
				Options(colorOptions()...).
				Value(color),
		),
	)

	return form.Run()
}

func ValidateColor(color string) bool {
	for _, c := range render.Names() {
		if c == color {
			return true
		}
	}
	return false
}
