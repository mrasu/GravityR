package db

import (
	"github.com/mrasu/GravityR/cmd/db/dig"
	"github.com/spf13/cobra"
)

var DigCmd = &cobra.Command{
	Use:   "dig",
	Short: "Dig database behavior",
}

func init() {
	DigCmd.AddCommand(dig.PerformanceInsightsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// SuggestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// flg := DigCmd.Flags()
	// flg.BoolVar(&runsExamination, "with-examine", false, "Examine query by adding index")
}
