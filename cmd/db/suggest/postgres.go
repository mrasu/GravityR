package suggest

import (
	"fmt"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/postgres"
	"github.com/mrasu/GravityR/database/postgres/model"
	"github.com/mrasu/GravityR/database/postgres/model/collector"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	iPostgres "github.com/mrasu/GravityR/infra/postgres"
	"github.com/spf13/cobra"
	"os"
	"path"
)

type postgresVarS struct {
	runsExamination bool
	indexTargets    []string
	query           string
}

var PostgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Suggest ways to increase PostgreSQL's performance",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := postgresRunner{}
		return r.run()
	},
}

var postgresVar = postgresVarS{}

func init() {
	flg := PostgresCmd.Flags()
	flg.BoolVar(&postgresVar.runsExamination, "with-examine", false, "Examine query by adding index")
	flg.StringArrayVarP(&postgresVar.indexTargets, "index", "i", []string{}, "Specify index")

	flg.StringVarP(&postgresVar.query, "query", "q", "", "[Required] Query to check")
	err := cobra.MarkFlagRequired(flg, "query")
	if err != nil {
		panic(err)
	}
}

type postgresRunner struct{}

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

	return pr.runSuggest(postgresVar, flag.DbFlag.Output, db, cfg.GetSearchPathOrPublic())
}

func (pr *postgresRunner) runSuggest(v postgresVarS, outputPath string, db *iPostgres.DB, schema string) error {
	examinationIdxTargets, err := parseIndexTargets(v.indexTargets)
	if err != nil {
		return err
	}

	explainLines, err := db.ExplainWithAnalyze(v.query)
	if err != nil {
		return err
	}

	aTree, err := collector.CollectExplainAnalyzeTree(explainLines)
	if err != nil {
		return err
	}

	its, errs := postgres.SuggestIndex(db, schema, v.query, aTree)
	if len(errs) > 0 {
		return errs[0]
	}
	idxTargets := toUniqueIndexTargets(its)

	fmt.Println("======suggest index by order-------")
	if len(idxTargets) > 0 {
		for _, it := range idxTargets {
			fmt.Println(it)
		}
	} else {
		fmt.Println("No suggestion. Perhaps already indexed?")
	}

	if len(examinationIdxTargets) == 0 {
		for _, it := range idxTargets {
			if it.IsSafe() {
				examinationIdxTargets = append(examinationIdxTargets, it)
			}
		}
	}

	var er *common_model.ExaminationResult
	if v.runsExamination {
		fmt.Printf("\n======going to examine-------\n")
		ie := postgres.NewIndexExaminer(db, v.query)
		er, err = database.NewIndexEfficiencyExaminer(ie).Run(examinationIdxTargets)
		if err != nil {
			return err
		}
	}

	if outputPath != "" {
		err := pr.createHTML(v, outputPath, idxTargets, er, aTree)
		if err != nil {
			return err
		}

		wd, err := os.Getwd()
		if err == nil {
			fmt.Printf("Result html is at: %s\n", path.Join(wd, outputPath))
		}
	}
	return nil
}

func (pr *postgresRunner) createHTML(v postgresVarS, outputPath string, idxTargets []*common_model.IndexTarget, er *common_model.ExaminationResult, aTree *model.ExplainAnalyzeTree) error {
	var vits []*viewmodel.VmIndexTarget
	for _, it := range idxTargets {
		vits = append(vits, it.ToViewModel())
	}

	var ver *viewmodel.VmExaminationResult
	if er != nil {
		ver = er.ToViewModel()
	}

	bo := html.NewSuggestPostgresDataBuildOption(
		v.query,
		aTree.ToViewModel(),
		aTree.SummaryText,
		vits,
		[]*viewmodel.VmExaminationCommandOption{
			viewmodel.CreateOutputExaminationOption(!v.runsExamination, outputPath),
			{IsShort: true, Name: "q", Value: v.query},
		},
		ver,
	)

	err := html.CreateHtml(outputPath, bo)
	if err != nil {
		return err
	}

	return nil
}
