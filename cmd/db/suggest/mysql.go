package suggest

import (
	"fmt"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/mysql"
	"github.com/mrasu/GravityR/database/mysql/model/collector"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	iMysql "github.com/mrasu/GravityR/infra/mysql"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var MySqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Suggest ways to increase MySQL's performance",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mysqlR.run()
	},
}

func init() {
	flg := MySqlCmd.Flags()
	flg.BoolVar(&mysqlR.runsExamination, "with-examine", false, "Examine query by adding index")
	flg.StringArrayVarP(&mysqlR.indexTargets, "index", "i", []string{}, "Specify index")

	flg.StringVarP(&mysqlR.query, "query", "q", "", "[Required] Query to check")
	err := cobra.MarkFlagRequired(flg, "query")
	if err != nil {
		panic(err)
	}
}

var mysqlR = mysqlRunner{}

type mysqlRunner struct {
	runsExamination bool
	indexTargets    []string
	query           string
}

func (mr *mysqlRunner) run() error {
	cfg, err := iMysql.NewConfigFromEnv()
	if err != nil {
		return err
	}

	db, err := iMysql.OpenMySQLDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	return mr.suggest(flag.DbFlag.Output, db, cfg.GetDBName())
}

func (mr *mysqlRunner) suggest(outputPath string, db *iMysql.DB, dbName string) error {
	examinationIdxTargets, err := parseIndexTargets(mr.indexTargets)
	if err != nil {
		return err
	}

	explainLine, err := db.Explain(mr.query)
	if err != nil {
		return err
	}

	aTree, err := collector.CollectExplainAnalyzeTree(explainLine)
	if err != nil {
		return err
	}

	its, errs := mysql.SuggestIndex(db, dbName, mr.query, aTree)
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
	if mr.runsExamination {
		fmt.Printf("\n======going to examine-------\n")
		ie := mysql.NewIndexExaminer(db, mr.query)
		er, err = database.NewIndexEfficiencyExaminer(ie).Run(examinationIdxTargets)
		if err != nil {
			return err
		}
	}

	if outputPath != "" {
		var vits []*viewmodel.VmIndexTarget
		for _, it := range idxTargets {
			vits = append(vits, it.ToViewModel())
		}

		var ver *viewmodel.VmExaminationResult
		if er != nil {
			ver = er.ToViewModel()
		}

		bo := html.NewSuggestMySQLDataBuildOption(
			mr.query,
			aTree.ToViewModel(),
			vits,
			[]*viewmodel.VmExaminationCommandOption{
				viewmodel.CreateOutputExaminationOption(!mr.runsExamination, outputPath),
				{IsShort: true, Name: "q", Value: mr.query},
			},
			ver,
		)

		err = html.CreateHtml(outputPath, bo)
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
