package suggest

import (
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/cmd/util"
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/database/postgres/pmodel"
	"github.com/mrasu/GravityR/database/postgres/pservice"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/mrasu/GravityR/infra/postgres"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
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
	cfg, err := postgres.NewConfigFromEnv()
	if err != nil {
		return err
	}

	db, err := postgres.OpenPostgresDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	return pr.suggest(flag.DbFlag.Output, db, cfg.GetSearchPathOrPublic())
}

func (pr *postgresRunner) suggest(outputPath string, db *postgres.DB, schema string) error {
	examinationIdxTargets, err := dmodel.NewIndexTargetsFromTexts(pr.indexTargets)
	if err != nil {
		return err
	}

	aTree, err := pservice.NewExplainer(db).ExplainWithAnalyze(pr.query)
	if err != nil {
		return err
	}

	its, err := pservice.NewIndexSuggester(db, schema).Suggest(pr.query)
	if err != nil {
		return err
	}
	util.LogNewIndexTargets(its)

	var er *dmodel.ExaminationResult
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

		util.LogResultOutputPath(outputPath)
	}
	return nil
}

func (pr *postgresRunner) examine(db *postgres.DB, varTargets, possibleTargets []*dmodel.IndexTarget) (*dmodel.ExaminationResult, error) {
	targets := varTargets
	if len(targets) == 0 {
		targets = lo.Filter(possibleTargets, func(it *dmodel.IndexTarget, _ int) bool { return it.IsSafe() })
	}

	log.Info().Msg("Start examination...")
	ie := pservice.NewIndexExaminer(db, pr.query)
	er, err := dservice.NewIndexEfficiencyExaminer(ie).Run(targets)
	if err != nil {
		return nil, err
	}

	return er, nil
}

func (pr *postgresRunner) createHTML(outputPath string, idxTargets []*dmodel.IndexTarget, er *dmodel.ExaminationResult, aTree *pmodel.ExplainAnalyzeTree) error {
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
