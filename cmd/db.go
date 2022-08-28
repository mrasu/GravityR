package cmd

import (
	"github.com/mrasu/GravityR/cmd/db"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Inspect DB-ish tool",
}

func init() {
	dbCmd.AddCommand(db.SuggestCmd)
	dbCmd.AddCommand(db.DigCmd)

	flg := dbCmd.PersistentFlags()
	flg.StringVarP(&flag.DbFlag.Output, "output", "o", "", "[Required] File name to output result html")
	cobra.MarkFlagRequired(flg, "output")
}
