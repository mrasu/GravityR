package suggest

import (
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/cmd/util"
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/mysql"
	"github.com/mrasu/GravityR/database/mysql/model/collector"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	iMysql "github.com/mrasu/GravityR/infra/mysql"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var MySqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Suggest ways to increase MySQL's performance",
	Run: func(cmd *cobra.Command, args []string) {
		err := mysqlR.run()
		if err != nil {
			util.LogError(err)
		}
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

	itts, errs := mysql.SuggestIndex(db, dbName, mr.query, aTree)
	if len(errs) > 0 {
		return errs[0]
	}

	its, err := mr.removeExistingIndexTargets(db, dbName, itts)
	if err != nil {
		return err
	}

	var er *common_model.ExaminationResult
	if mr.runsExamination {
		er, err = mr.examine(db, examinationIdxTargets, its)
		if err != nil {
			return err
		}
	}

	if outputPath != "" {
		vits := lo.Map(its, func(it *common_model.IndexTarget, _ int) *viewmodel.VmIndexTarget { return it.ToViewModel() })

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
			log.Info().Msg("Result html is at: " + path.Join(wd, outputPath))
		}
	}
	return nil
}

func (mr *mysqlRunner) removeExistingIndexTargets(db *iMysql.DB, dbName string, itts []*common_model.IndexTargetTable) ([]*common_model.IndexTarget, error) {
	idxGetter := mysql.NewIndexGetter(db)
	its, err := database.NewExistingIndexRemover(idxGetter, dbName, itts).Remove()
	if err != nil {
		return nil, err
	}

	logNewIndexTargets(its)
	return its, nil
}

func (mr *mysqlRunner) examine(db *iMysql.DB, varTargets, possibleTargets []*common_model.IndexTarget) (*common_model.ExaminationResult, error) {
	targets := varTargets
	if len(targets) == 0 {
		targets = lo.Filter(possibleTargets, func(it *common_model.IndexTarget, _ int) bool { return it.IsSafe() })
	}

	log.Info().Msg("Start examination...")
	ie := mysql.NewIndexExaminer(db, mr.query)
	er, err := database.NewIndexEfficiencyExaminer(ie).Run(targets)
	if err != nil {
		return nil, err
	}

	return er, nil
}
