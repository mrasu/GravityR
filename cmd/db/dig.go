package db

import (
	"github.com/mrasu/GravityR/cmd/db/dig"
	"github.com/spf13/cobra"
)

// DigCmd represents the db command
var DigCmd = &cobra.Command{
	Use:   "dig",
	Short: "Dig database behavior",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is IAnalyzeData CLI library for Go that empowers applications.
This application is IAnalyzeData tool to generate the needed files
to quickly create IAnalyzeData Cobra application.`,
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
