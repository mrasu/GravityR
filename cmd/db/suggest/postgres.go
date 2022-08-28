package suggest

import (
	"fmt"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/database/postgres"
	"github.com/mrasu/GravityR/database/postgres/models"
	"github.com/mrasu/GravityR/database/postgres/models/collectors"
	"github.com/mrasu/GravityR/database/rdb"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
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
	cobra.MarkFlagRequired(flg, "query")
}

type postgresRunner struct{}

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

	examinationIdxTargets, err := pr.parseIndexTargets(postgresVar.indexTargets)
	if err != nil {
		return err
	}

	aTree, err := collectors.CollectExplainAnalyzeTree(db, postgresVar.query)
	if err != nil {
		return err
	}

	its, errs := postgres.SuggestIndex(db, cfg.GetSearchPathOrPublic(), postgresVar.query, aTree)
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

	var er *db_models.ExaminationResult
	if postgresVar.runsExamination {
		fmt.Printf("\n======going to examine-------\n")
		ie := &rdb.IndexExaminer{DbType: rdb.PostgreSQL}
		er, err = ie.Run(db, postgresVar.query, examinationIdxTargets)
		if err != nil {
			return err
		}
	}

	outputPath := flag.DbFlag.Output
	if outputPath != "" {
		err := pr.createHTML(outputPath, idxTargets, er, aTree)
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

func (pr *postgresRunner) parseIndexTargets(indexTargetTexts []string) ([]*db_models.IndexTarget, error) {
	var its []*db_models.IndexTarget
	for _, text := range indexTargetTexts {
		it, err := db_models.NewIndexTarget(text)
		if err != nil {
			return nil, err
		}
		its = append(its, it)
	}

	return its, nil
}

func (pr *postgresRunner) createHTML(outputPath string, idxTargets []*db_models.IndexTarget, er *db_models.ExaminationResult, aTree *models.ExplainAnalyzeTree) error {
	var vits []*viewmodel.VmIndexTarget
	for _, it := range idxTargets {
		vits = append(vits, it.ToViewModel())
	}

	var ver *viewmodel.VmExaminationResult
	if er != nil {
		ver = er.ToViewModel()
	}

	bo := html.NewSuggestPostgresDataBuildOption(
		postgresVar.query,
		aTree.ToViewModel(),
		aTree.PlanningText,
		vits,
		[]*viewmodel.VmExaminationCommandOption{
			viewmodel.CreateOutputExaminationOption(!postgresVar.runsExamination, outputPath),
			{IsShort: true, Name: "q", Value: postgresVar.query},
		},
		ver,
	)

	err := html.CreateHtml(outputPath, bo)
	if err != nil {
		return err
	}

	return nil
}
