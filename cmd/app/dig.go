package app

import (
	"github.com/mrasu/GravityR/cmd/app/dig"
	"github.com/spf13/cobra"
)

var DigCmd = &cobra.Command{
	Use:   "dig",
	Short: "Dig application behavior",
}

func init() {
	DigCmd.AddCommand(dig.JaegerCmd)
}
