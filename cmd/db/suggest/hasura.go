package suggest

import (
	"github.com/mrasu/GravityR/cmd/db/suggest/hasura"
	"github.com/spf13/cobra"
)

var HasuraCmd = &cobra.Command{
	Use:   "hasura",
	Short: "Suggest ways to increase Hasura's performance",
}

func init() {
	HasuraCmd.AddCommand(hasura.PostgresCmd)
}
