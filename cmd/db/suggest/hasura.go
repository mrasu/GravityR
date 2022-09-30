package suggest

import (
	"encoding/json"
	"fmt"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/database"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/hasura"
	"github.com/mrasu/GravityR/database/postgres/model"
	"github.com/mrasu/GravityR/database/postgres/model/collector"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	iHasura "github.com/mrasu/GravityR/infra/hasura"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var HasuraCmd = &cobra.Command{
	Use:   "hasura",
	Short: "Suggest ways to increase Hasura's performance",
	RunE: func(cmd *cobra.Command, args []string) error {
		return hasuraR.run()
	},
}

func init() {
	flg := HasuraCmd.Flags()
	flg.BoolVar(&hasuraR.runsExamination, "with-examine", false, "Examine query by adding index")
	flg.StringArrayVarP(&hasuraR.indexTargets, "index", "i", []string{}, "Specify index")

	flg.StringVarP(&hasuraR.query, "query", "q", "", "[Required] Query of GraphQL to check")
	flg.StringVar(&hasuraR.jsonVariables, "json-variables", "{}", "JSON of variables for the query")
	err := cobra.MarkFlagRequired(flg, "query")
	if err != nil {
		panic(err)
	}
}

var hasuraR = hasuraRunner{}

type hasuraRunner struct {
	runsExamination bool
	indexTargets    []string
	query           string
	jsonVariables   string
}

func (hr *hasuraRunner) run() error {
	cfg, err := iHasura.NewConfigFromEnv()
	if err != nil {
		return err
	}
	cli := iHasura.NewClient(cfg)

	return hr.suggest(flag.DbFlag.Output, cli)
}

func (hr *hasuraRunner) suggest(outputPath string, cli *iHasura.Client) error {
	v, err := hr.parseJSONToVariables(hr.jsonVariables)
	if err != nil {
		return err
	}

	q := &iHasura.ExplainRequestBody{
		Query: &iHasura.Query{
			Query:     hr.query,
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
	aTree, err := collector.CollectExplainAnalyzeTree(r.Plan)
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

	examinationIdxTargets, err := parseIndexTargets(hr.indexTargets)
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

	var er *common_model.ExaminationResult
	if hr.runsExamination {
		fmt.Printf("\n======going to examine-------\n")
		ie := hasura.NewIndexExaminer(cli, hr.query, v)
		er, err = database.NewIndexEfficiencyExaminer(ie).Run(examinationIdxTargets)
		if err != nil {
			return err
		}
	}

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
		return nil, errors.Wrap(err, "failed to parse variables")
	}

	return variables, nil
}

func (hr *hasuraRunner) createHTML(outputPath string, variables map[string]interface{}, sql string, idxTargets []*common_model.IndexTarget, er *common_model.ExaminationResult, aTree *model.ExplainAnalyzeTree) error {
	vits := lo.Map(idxTargets, func(v *common_model.IndexTarget, _ int) *viewmodel.VmIndexTarget { return v.ToViewModel() })

	var ver *viewmodel.VmExaminationResult
	if er != nil {
		ver = er.ToViewModel()
	}

	bo := html.NewSuggestHasuraDataBuildOption(
		hr.query,
		variables,
		sql,
		aTree.ToViewModel(),
		aTree.SummaryText,
		vits,
		[]*viewmodel.VmExaminationCommandOption{
			viewmodel.CreateOutputExaminationOption(!hr.runsExamination, outputPath),
			{IsShort: true, Name: "q", Value: hr.query},
			{IsShort: false, Name: "json-variables", Value: hr.jsonVariables},
		},
		ver,
	)

	err := html.CreateHtml(outputPath, bo)
	if err != nil {
		return err
	}

	return nil
}
