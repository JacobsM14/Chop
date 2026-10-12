package cmd

import (
	"strconv"

	"github.com/JacobsM14/chop/internal/store"

	"github.com/spf13/cobra"
)

func completeChopNames(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	chops, err := db.List(store.SortByName)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	names := make([]cobra.Completion, len(chops))
	for i, chop := range chops {
		names[i] = chop.Name
	}
	return names, cobra.ShellCompDirectiveNoFileComp
}

func completeItemNumbers(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	chop, err := db.Get(args[0])
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	items := make([]cobra.Completion, len(chop.Items))
	for i, it := range chop.Items {
		items[i] = cobra.CompletionWithDesc(strconv.Itoa(i+1), it.Text)
	}
	return items, cobra.ShellCompDirectiveNoFileComp
}
