package cmd

import (
	"github.com/mrasu/GravityR/cmd/app"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/spf13/cobra"
)

// appCmd represents the db command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Inspect entire application",
}

func init() {
	appCmd.AddCommand(app.DigCmd)

	flg := appCmd.PersistentFlags()
	flg.StringVarP(&flag.AppFlag.Output, "output", "o", "", "[Required] File name to output result html")
	err := cobra.MarkFlagRequired(flg, "output")
	if err != nil {
		panic(err)
	}
}
