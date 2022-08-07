package dig

import (
	"fmt"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/infra/aws"
	"github.com/mrasu/GravityR/infra/aws/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"time"
)

var PerformanceInsightsCmd = &cobra.Command{
	Use:   "performance-insights",
	Short: "Dig database behavior with AWS' PerformanceInsights",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is IAnalyzeData CLI library for Go that empowers applications.
This application is IAnalyzeData tool to generate the needed files
to quickly create IAnalyzeData Cobra application.`,
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
	fmt.Println(flag.GlobalFlag.GetAwsEndpoint())

	gFlg := flag.GlobalFlag
	cfg, err := aws.NewAwsConfig(gFlg.GetAwsEndpoint(), gFlg.Verbose)
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
	var sqlLoadAvgs []*models.PiSQLLoadAvg
	var tokenizedSqlLoadAvgs []*models.PiSQLLoadAvg
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

	sqlDbLoads := models.ConvertPiSQLLoadAvgsToVms(startFrom, end, sqlLoadAvgs)
	tokenizedSqlDbLoads := models.ConvertPiSQLLoadAvgsToVms(startFrom, end, tokenizedSqlLoadAvgs)
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
