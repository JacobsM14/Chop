package cmd

import (
	"errors"
	"fmt"

	"github.com/JacobsM14/chop/internal/render"
	"github.com/JacobsM14/chop/internal/store"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:               "get <name>",
	Aliases:           []string{"ver"},
	Short:             "Show a Chop and its items",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeChopNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		chop, err := db.Get(args[0])
		if errors.Is(err, store.ErrNotFound) {
			return fmt.Errorf("chop %q not found", args[0])
		}
		if err != nil {
			return err
		}

		fmt.Println(render.Title(chop.Name, chop.MainColor))
		if len(chop.Items) == 0 {
			fmt.Printf("no items yet, add one with: chop add %s\n", chop.Name)
			return nil
		}
		for i, it := range chop.Items {
			fmt.Println(render.Item(i, chop.Ordered, it))
		}
		return nil
	},
}
