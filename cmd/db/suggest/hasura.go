package suggest

import (
	"encoding/json"
	"fmt"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/database/db_models"
	"github.com/mrasu/GravityR/database/hasura"
	"github.com/mrasu/GravityR/database/postgres/models"
	"github.com/mrasu/GravityR/database/postgres/models/collectors"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	iHasura "github.com/mrasu/GravityR/infra/hasura"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var HasuraCmd = &cobra.Command{
	Use:   "hasura",
	Short: "Suggest ways to increase Hasura's performance",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := hasuraRunner{}
		return r.run()
	},
}

type hasuraVarS struct {
	runsExamination bool
	indexTargets    []string
	query           string
	jsonVariables   string
}

var hasuraVar = hasuraVarS{}

func init() {
	flg := HasuraCmd.Flags()
	flg.BoolVar(&hasuraVar.runsExamination, "with-examine", false, "Examine query by adding index")
	flg.StringArrayVarP(&hasuraVar.indexTargets, "index", "i", []string{}, "Specify index")

	flg.StringVarP(&hasuraVar.query, "query", "q", "", "[Required] Query of GraphQL to check")
	flg.StringVar(&hasuraVar.jsonVariables, "json-variables", "{}", "JSON of variables for the query")
	cobra.MarkFlagRequired(flg, "query")
}

type hasuraRunner struct{}

func (hr *hasuraRunner) run() error {
	cfg, err := iHasura.NewConfigFromEnv()
	if err != nil {
		return err
	}
	cli := iHasura.NewClient(cfg)

	v, err := hr.parseJSONToVariables(hasuraVar.jsonVariables)
	if err != nil {
		return err
	}

	q := &iHasura.ExplainRequestBody{
		Query: &iHasura.Query{
			Query:     hasuraVar.query,
			Variables: v,
		},
	}
	res, err := cli.Explain(q)
	if err != nil {
		return err
	}
	if len(res) != 1 {
		return errors.Errorf("not supporting multiple gql executing multiple queries. %d query will be run", len(res))
	}

	r := res[0]
	aTree, err := collectors.CollectExplainAnalyzeTree(r.Plan)
	if err != nil {
		return err
	}

	fmt.Println("======explain-------")
	for _, p := range r.Plan {
		fmt.Println(p)
	}

	its, errs := hasura.SuggestIndex(cli, r.SQL, aTree)
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

	examinationIdxTargets, err := parseIndexTargets(hasuraVar.indexTargets)
	if err != nil {
		return err
	}

	if len(examinationIdxTargets) == 0 {
		for _, it := range idxTargets {
			if it.IsSafe() {
				examinationIdxTargets = append(examinationIdxTargets, it)
			}
		}
	}

	var er *db_models.ExaminationResult
	if hasuraVar.runsExamination {
		fmt.Printf("\n======going to examine-------\n")
		ie := hasura.NewIndexExaminer(cli, hasuraVar.query, v)
		er, err = database.NewIndexEfficiencyExaminer(ie).Run(examinationIdxTargets)
		if err != nil {
			return err
		}
	}

	outputPath := flag.DbFlag.Output
	if outputPath != "" {
		err := hr.createHTML(outputPath, v, r.SQL, idxTargets, er, aTree)
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

func (hr *hasuraRunner) parseJSONToVariables(jsonStr string) (map[string]interface{}, error) {
	var variables map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &variables); err != nil {
		return nil, err
	}

	return variables, nil
}

func (hr *hasuraRunner) createHTML(outputPath string, variables map[string]interface{}, sql string, idxTargets []*db_models.IndexTarget, er *db_models.ExaminationResult, aTree *models.ExplainAnalyzeTree) error {
	vits := lib.Map(idxTargets, func(v *db_models.IndexTarget) *viewmodel.VmIndexTarget { return v.ToViewModel() })

	var ver *viewmodel.VmExaminationResult
	if er != nil {
		ver = er.ToViewModel()
	}

	bo := html.NewSuggestHasuraDataBuildOption(
		hasuraVar.query,
		variables,
		sql,
		aTree.ToViewModel(),
		aTree.SummaryText,
		vits,
		[]*viewmodel.VmExaminationCommandOption{
			viewmodel.CreateOutputExaminationOption(!hasuraVar.runsExamination, outputPath),
			{IsShort: true, Name: "q", Value: hasuraVar.query},
			{IsShort: false, Name: "json-variables", Value: hasuraVar.jsonVariables},
		},
		ver,
	)

	err := html.CreateHtml(outputPath, bo)
	if err != nil {
		return err
	}

	return nil
}