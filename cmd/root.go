package cmd

import (
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/cmd/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "gr",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Gravity Radar to remove bottleneck in your application",
	Long: `GravityR is Gravity-Radar.
This exists to reduce time to find bottleneck in your application.
And also this is to solve the problems faster and easier.  
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		if flag.GlobalFlag.Verbose {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		if flag.GlobalFlag.UseMock {
			log.Info().Msg("NOTE: Using mock")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		util.LogError(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().BoolVar(&flag.GlobalFlag.UseMock, "use-mock", false, "use mock (for development)")
	rootCmd.PersistentFlags().BoolVarP(&flag.GlobalFlag.Verbose, "verbose", "v", false, "verbose output")
}
