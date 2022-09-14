package db

import (
	"github.com/mrasu/GravityR/cmd/db/suggest"
	"github.com/spf13/cobra"
)

var SuggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Inspect and suggest",
}

func init() {
	SuggestCmd.AddCommand(suggest.MySqlCmd)
	SuggestCmd.AddCommand(suggest.PostgresCmd)
	SuggestCmd.AddCommand(suggest.HasuraCmd)
}
