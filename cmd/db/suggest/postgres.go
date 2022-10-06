package suggest

import (
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/cmd/util"
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/database/postgres"
	"github.com/mrasu/GravityR/database/postgres/model"
	"github.com/mrasu/GravityR/database/postgres/model/collector"
	"github.com/mrasu/GravityR/database/service"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	iPostgres "github.com/mrasu/GravityR/infra/postgres"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var PostgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Suggest ways to increase PostgreSQL's performance",
	Run: func(cmd *cobra.Command, args []string) {
		err := postgresR.run()
		if err != nil {
			util.LogError(err)
		}
	},
}

func init() {
	flg := PostgresCmd.Flags()
	flg.BoolVar(&postgresR.runsExamination, "with-examine", false, "Examine query by adding index")
	flg.StringArrayVarP(&postgresR.indexTargets, "index", "i", []string{}, "Specify index")

	flg.StringVarP(&postgresR.query, "query", "q", "", "[Required] Query to check")
	err := cobra.MarkFlagRequired(flg, "query")
	if err != nil {
		panic(err)
	}
}

var postgresR = postgresRunner{}

type postgresRunner struct {
	runsExamination bool
	indexTargets    []string
	query           string
}

func (pr *postgresRunner) run() error {
	cfg, err := iPostgres.NewConfigFromEnv()
	if err != nil {
		return err
	}

	db, err := iPostgres.OpenPostgresDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	return pr.suggest(flag.DbFlag.Output, db, cfg.GetSearchPathOrPublic())
}

func (pr *postgresRunner) suggest(outputPath string, db *iPostgres.DB, schema string) error {
	examinationIdxTargets, err := parseIndexTargets(pr.indexTargets)
	if err != nil {
		return err
	}

	explainLines, err := db.ExplainWithAnalyze(pr.query)
	if err != nil {
		return err
	}

	aTree, err := collector.CollectExplainAnalyzeTree(explainLines)
	if err != nil {
		return err
	}

	itts, errs := postgres.SuggestIndex(db, schema, pr.query, aTree)
	if len(errs) > 0 {
		return errs[0]
	}

	its, err := pr.removeExistingIndexTargets(db, schema, itts)
	if err != nil {
		return err
	}

	var er *database.ExaminationResult
	if pr.runsExamination {
		er, err = pr.examine(db, examinationIdxTargets, its)
		if err != nil {
			return err
		}
	}

	if outputPath != "" {
		err := pr.createHTML(outputPath, its, er, aTree)
		if err != nil {
			return err
		}

		wd, err := os.Getwd()
		if err == nil {
			log.Info().Msg("Result html is at: " + path.Join(wd, outputPath))
		}
	}
	return nil
}

func (pr *postgresRunner) removeExistingIndexTargets(db *iPostgres.DB, dbName string, itts []*database.IndexTargetTable) ([]*database.IndexTarget, error) {
	idxGetter := postgres.NewIndexGetter(db)
	its, err := service.NewExistingIndexRemover(idxGetter, dbName, itts).Remove()
	if err != nil {
		return nil, err
	}

	logNewIndexTargets(its)
	return its, nil
}

func (pr *postgresRunner) examine(db *iPostgres.DB, varTargets, possibleTargets []*database.IndexTarget) (*database.ExaminationResult, error) {
	targets := varTargets
	if len(targets) == 0 {
		targets = lo.Filter(possibleTargets, func(it *database.IndexTarget, _ int) bool { return it.IsSafe() })
	}

	log.Info().Msg("Start examination...")
	ie := postgres.NewIndexExaminer(db, pr.query)
	er, err := service.NewIndexEfficiencyExaminer(ie).Run(targets)
	if err != nil {
		return nil, err
	}

	return er, nil
}

func (pr *postgresRunner) createHTML(outputPath string, idxTargets []*database.IndexTarget, er *database.ExaminationResult, aTree *model.ExplainAnalyzeTree) error {
	var vits []*viewmodel.VmIndexTarget
	for _, it := range idxTargets {
		vits = append(vits, it.ToViewModel())
	}

	var ver *viewmodel.VmExaminationResult
	if er != nil {
		ver = er.ToViewModel()
	}

	bo := html.NewSuggestPostgresDataBuildOption(
		pr.query,
		aTree.ToViewModel(),
		aTree.SummaryText,
		vits,
		[]*viewmodel.VmExaminationCommandOption{
			viewmodel.CreateOutputExaminationOption(!pr.runsExamination, outputPath),
			{IsShort: true, Name: "q", Value: pr.query},
		},
		ver,
	)

	err := html.CreateHtml(outputPath, bo)
	if err != nil {
		return err
	}

	return nil
}
