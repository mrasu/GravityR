package suggest

import (
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/cmd/lib"
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/postgres"
	"github.com/mrasu/GravityR/database/postgres/model"
	"github.com/mrasu/GravityR/database/postgres/model/collector"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	iPostgres "github.com/mrasu/GravityR/infra/postgres"
	"github.com/rs/zerolog/log"
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
			lib.LogError(err)
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

	its, errs := postgres.SuggestIndex(db, schema, pr.query, aTree)
	if len(errs) > 0 {
		return errs[0]
	}
	idxTargets := toUniqueIndexTargets(its)

	if len(idxTargets) > 0 {
		log.Debug().Msg("Found possibly efficient index combinations:")
		for i, it := range idxTargets {
			log.Printf("\t%d.%s", i, it.CombinationString())
		}
	} else {
		log.Debug().Msg("No possibly efficient index found. Perhaps already indexed?")
	}

	if len(examinationIdxTargets) == 0 {
		for _, it := range idxTargets {
			if it.IsSafe() {
				examinationIdxTargets = append(examinationIdxTargets, it)
			}
		}
	}

	var er *common_model.ExaminationResult
	if pr.runsExamination {
		log.Info().Msg("Start examination...")
		ie := postgres.NewIndexExaminer(db, pr.query)
		er, err = database.NewIndexEfficiencyExaminer(ie).Run(examinationIdxTargets)
		if err != nil {
			return err
		}
	}

	if outputPath != "" {
		err := pr.createHTML(outputPath, idxTargets, er, aTree)
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

func (pr *postgresRunner) createHTML(outputPath string, idxTargets []*common_model.IndexTarget, er *common_model.ExaminationResult, aTree *model.ExplainAnalyzeTree) error {
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
