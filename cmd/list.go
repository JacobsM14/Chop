package cmd

import (
	"fmt"

	"github.com/JacobsM14/chop/internal/render"
	"github.com/JacobsM14/chop/internal/store"
	"github.com/spf13/cobra"
)

var sortFlag string

func init() {
	lsCmd.Flags().StringVar(&sortFlag, "sort", "", "sort by: name, color, ordered")
	sorts := []cobra.Completion{string(store.SortByName), string(store.SortByColor), string(store.SortByOrdered)}
	cobra.CheckErr(lsCmd.RegisterFlagCompletionFunc("sort", cobra.FixedCompletions(sorts, cobra.ShellCompDirectiveNoFileComp)))
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:               "ls",
	Short:             "List all Chops",
	ValidArgsFunction: cobra.NoFileCompletions,
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
