package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "pkgrepo",
	Short: "Image Packaging System repository management utility",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pkgrepo: no subcommand specified")
		fmt.Println("Try `pkgrepo --help or -?' for more information.")
		os.Exit(2)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
