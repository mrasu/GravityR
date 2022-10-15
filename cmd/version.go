package cmd

import (
	"fmt"
	"github.com/mrasu/GravityR/injection"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GravityR v%s, commit %s\n", injection.BinInfo.Version, injection.BinInfo.Commit)
	},
}
