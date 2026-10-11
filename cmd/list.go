package cmd

import (
	"fmt"

	"chop/internal/render"
	"chop/internal/store"

	"github.com/spf13/cobra"
)

var sortFlag string

func init() {
	lsCmd.Flags().StringVar(&sortFlag, "sort", "", "sort by: name, color, ordered")
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all Chops",
	RunE: func(cmd *cobra.Command, args []string) error {
		chops, err := db.List(store.SortField(sortFlag))
		if err != nil {
			return err
		}
		for _, chop := range chops {
			fmt.Println(render.Paint(chop.MainColor, chop.Name))
		}
		return nil
	},
}
