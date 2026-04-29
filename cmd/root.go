package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vsc",
	Short: "vsc is a CLI tool for VS Code",
	Long:  `vsc — interact with VS Code from the command line`,
}

func init() {
	rootCmd.AddCommand(lsCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}