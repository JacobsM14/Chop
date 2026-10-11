package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	newColor   string
	newOrdered bool
)

var colors = []string{"red", "green", "blue", "purple", "orange", "white"}

var colorCodes = map[string]string{
	"red": "9", "green": "10", "blue": "12",
	"purple": "13", "orange": "208", "white": "15",
}

// Initialize the new command and its flags
func init() {
	newCmd.Flags().StringVar(&newColor, "color", "", "color principal")
	newCmd.Flags().BoolVar(&newOrdered, "ordered", false, "lista numerada")
	rootCmd.AddCommand(newCmd)
}

var newCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create a new Chop",
	Args:  cobra.ExactArgs(1),
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

		if !ValidateColor(newColor) {
			return fmt.Errorf("invalid color %q (valid: %s)", newColor, strings.Join(colors, ", "))
		}

		fmt.Printf("name=%s color=%s ordered=%t\n", name, newColor, newOrdered)
		return nil
	},
}

func colorOptions() []huh.Option[string] {
	opts := make([]huh.Option[string], 0, len(colors))
	for _, c := range colors {
		label := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCodes[c])).
			Render(c)
		opts = append(opts, huh.NewOption(label, c))
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
	for _, c := range colors {
		if c == color {
			return true
		}
	}
	return false
}
