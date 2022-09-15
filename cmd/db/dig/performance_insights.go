package dig

import (
	"fmt"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/infra/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"time"
)

var PerformanceInsightsCmd = &cobra.Command{
	Use:   "performance-insights",
	Short: "Dig database behavior with AWS' PerformanceInsights",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDig()
	},
}

type piVarS struct {
	start string
}

var piVar = piVarS{}

func init() {
	PerformanceInsightsCmd.Flags().StringVar(&piVar.start, "start-from", "", "Date dig from (RFC3339) (Default: 7 days ago)")
}

func runDig() error {
	fmt.Printf("mocked?: %t\n", flag.GlobalFlag.UseMock)

	gFlg := flag.GlobalFlag
	cfg, err := aws.NewAwsConfig(gFlg.UseMock, gFlg.Verbose)
	if err != nil {
		return err
	}

	cli := aws.NewRds(cfg)
	dbs, err := cli.GetDBs([]string{"mysql"})
	if err != nil {
		return err
	}

	client := aws.NewPerformanceInsights(cfg)
	if err != nil {
		return err
	}

	startFrom := time.Now().Add(-7 * 24 * time.Hour)
	if piVar.start != "" {
		t, err := time.Parse(time.RFC3339, piVar.start)
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
			fmt.Printf("Getting data for %s\n", current.Format(time.RFC3339))

			avgs, err := client.GetHalfDaySqlMetrics(db, current)
			if err != nil {
				return err
			}
			sqlLoadAvgs = append(sqlLoadAvgs, avgs...)

			tAvgs, err := client.GetHalfDayTokenizedSqlMetrics(db, current)
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

	outputPath := flag.DbFlag.Output
	err = html.CreateHtml(outputPath, bo)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err == nil {
		fmt.Printf("Result html is at: %s\n", path.Join(wd, outputPath))
	}

	return nil
}
