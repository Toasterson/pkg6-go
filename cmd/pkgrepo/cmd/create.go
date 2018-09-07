package cmd

import (
	"git.wegmueller.it/toasterson/glog"
	"github.com/spf13/cobra"
	"github.com/toasterson/pkg6-go/repo"
	"os"
	"strings"
)

var createCMD = &cobra.Command{
	Use:   "create [--version ver] uri_or_path",
	Short: "Create a pkg(5) repository at the specified location.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pathArg := args[0]
		if strings.HasPrefix(pathArg, "file://") {
			pathArg = strings.Replace(pathArg, "file://", "", -1)
		}
		r, err := repo.NewRepo(pathArg)
		if err != nil {
			glog.Emergln(err)
			os.Exit(1)
		}
		if err = r.Create(); err != nil {
			glog.Errf("could not create repository at %s: %s", pathArg, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCMD)
}
