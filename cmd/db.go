package cmd

import (
	"github.com/mrasu/GravityR/cmd/db"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	dbCmd.AddCommand(db.SuggestCmd)
	dbCmd.AddCommand(db.SampleCmd)
	dbCmd.AddCommand(db.DigCmd)

	flg := dbCmd.PersistentFlags()
	flg.StringVarP(&flag.DbFlag.Output, "output", "o", "", "[Required] File name to output result html")
	cobra.MarkFlagRequired(flg, "output")
}
