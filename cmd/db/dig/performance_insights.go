package dig

import (
	"fmt"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/cmd/util"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/infra/aws"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"time"
)

var PerformanceInsightsCmd = &cobra.Command{
	Use:   "performance-insights",
	Short: "Dig database behavior with AWS' PerformanceInsights",
	RunE: func(cmd *cobra.Command, args []string) error {
		return piR.run()
	},
}

func init() {
	PerformanceInsightsCmd.Flags().StringVar(&piR.start, "start-from", "", "Date dig from (RFC3339) (Default: 7 days ago)")
}

var piR = piRunner{}

type piRunner struct {
	start string
}

func (pr *piRunner) run() error {
	gFlg := flag.GlobalFlag
	cfg, err := aws.NewAwsConfig(gFlg.UseMock, gFlg.Verbose)
	if err != nil {
		return err
	}

	rdsCli := aws.NewRds(cfg)
	piCli := aws.NewPerformanceInsights(cfg)

	return pr.dig(flag.DbFlag.Output, rdsCli, piCli)
}

func (pr *piRunner) dig(outputPath string, rdsCli *aws.Rds, piCli *aws.PerformanceInsights) error {
	dbs, err := rdsCli.GetDBs([]string{"mysql"})
	if err != nil {
		return err
	}

	startFrom := time.Now().Add(-7 * 24 * time.Hour)
	if pr.start != "" {
		t, err := time.Parse(time.RFC3339, pr.start)
		if err != nil {
			return errors.Wrap(err, "Invalid format")
		}
		startFrom = t
	}
	end := startFrom.Add(7 * 24 * time.Hour)
	var sqlLoadAvgs []*aws.PiSQLLoadAvg
	var tokenizedSqlLoadAvgs []*aws.PiSQLLoadAvg
	for _, db := range dbs {
		current := startFrom
		for i := 0; i < 14; i++ {
			log.Info().Msg(fmt.Sprintf("Getting data for %s", current.Format(time.RFC3339)))

			avgs, err := piCli.GetHalfDaySqlMetrics(db, current)
			if err != nil {
				return err
			}
			sqlLoadAvgs = append(sqlLoadAvgs, avgs...)

			tAvgs, err := piCli.GetHalfDayTokenizedSqlMetrics(db, current)
			if err != nil {
				return err
			}
			tokenizedSqlLoadAvgs = append(tokenizedSqlLoadAvgs, tAvgs...)

			current = current.Add(12 * time.Hour)
		}
	}

	sqlDbLoads := aws.ConvertPiSQLLoadAvgsToVms(startFrom, end, sqlLoadAvgs)
	tokenizedSqlDbLoads := aws.ConvertPiSQLLoadAvgsToVms(startFrom, end, tokenizedSqlLoadAvgs)
	bo := html.NewDigDataBuildOption(sqlDbLoads, tokenizedSqlDbLoads)

	err = html.CreateHtml(outputPath, bo)
	if err != nil {
		return err
	}

	util.LogResultOutputPath(outputPath)

	return nil
}
