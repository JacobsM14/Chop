package cmd

import (
	"chop/internal/store"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chop",
	Short: "...",
	// PersistentPreRunE: abre el Store
	// PersistentPostRunE: lo cierra
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// imprimir el error por stderr
		// salir con código 1
	}
}

func Open(path string) (store.Store, error) {
	return nil, nil
}
