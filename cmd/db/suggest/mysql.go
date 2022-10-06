package suggest

import (
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/cmd/util"
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/database/mysql/mservice"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/mrasu/GravityR/infra/mysql"
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
	cfg, err := mysql.NewConfigFromEnv()
	if err != nil {
		return err
	}

	db, err := mysql.OpenMySQLDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	return mr.suggest(flag.DbFlag.Output, db, cfg.GetDBName())
}

func (mr *mysqlRunner) suggest(outputPath string, db *mysql.DB, dbName string) error {
	examinationIdxTargets, err := parseIndexTargets(mr.indexTargets)
	if err != nil {
		return err
	}

	explainLine, err := db.Explain(mr.query)
	if err != nil {
		return err
	}

	aTree, err := mservice.CollectExplainAnalyzeTree(explainLine)
	if err != nil {
		return err
	}

	itts, errs := mservice.SuggestIndex(db, dbName, mr.query, aTree)
	if len(errs) > 0 {
		return errs[0]
	}

	its, err := mr.removeExistingIndexTargets(db, dbName, itts)
	if err != nil {
		return err
	}

	var er *dmodel.ExaminationResult
	if mr.runsExamination {
		er, err = mr.examine(db, examinationIdxTargets, its)
		if err != nil {
			return err
		}
	}

	if outputPath != "" {
		vits := lo.Map(its, func(it *dmodel.IndexTarget, _ int) *viewmodel.VmIndexTarget { return it.ToViewModel() })

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

func (mr *mysqlRunner) removeExistingIndexTargets(db *mysql.DB, dbName string, itts []*dmodel.IndexTargetTable) ([]*dmodel.IndexTarget, error) {
	idxGetter := mservice.NewIndexGetter(db)
	its, err := dservice.NewExistingIndexRemover(idxGetter, dbName, itts).Remove()
	if err != nil {
		return nil, err
	}

	logNewIndexTargets(its)
	return its, nil
}

func (mr *mysqlRunner) examine(db *mysql.DB, varTargets, possibleTargets []*dmodel.IndexTarget) (*dmodel.ExaminationResult, error) {
	targets := varTargets
	if len(targets) == 0 {
		targets = lo.Filter(possibleTargets, func(it *dmodel.IndexTarget, _ int) bool { return it.IsSafe() })
	}

	log.Info().Msg("Start examination...")
	ie := mservice.NewIndexExaminer(db, mr.query)
	er, err := dservice.NewIndexEfficiencyExaminer(ie).Run(targets)
	if err != nil {
		return nil, err
	}

	return er, nil
}
