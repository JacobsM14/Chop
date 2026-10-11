package cmd

import (
	"chop/internal/store"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		chopName := args[0]
		item, err := askItem()
		if errors.Is(err, huh.ErrUserAborted) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("running form: %w", err)
		}

		err = db.AddItem(chopName, item)
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("chop %q not found", chopName)
		}
		return err
	},
}

func askItem() (store.ListItem, error) {
	item := store.ListItem{TextColor: "purple", DescriptionColor: "white"}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Command").
				Value(&item.Text).Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return errors.New("command cannot be empty")
				}
				return nil
			}),

			huh.NewSelect[string]().
				Title("Text Color").
				Options(colorOptions()...).
				Value(&item.TextColor),

			huh.NewInput().
				Title("Description").
				Value(&item.Description).Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return errors.New("description cannot be empty")
				}
				return nil
			}),

			huh.NewSelect[string]().
				Title("Description Color").
				Options(colorOptions()...).
				Value(&item.DescriptionColor),
		),
	)

	if err := form.Run(); err != nil {
		return store.ListItem{}, err
	}
	return item, nil
}
