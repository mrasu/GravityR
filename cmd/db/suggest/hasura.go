package suggest

import (
	"encoding/json"
	"github.com/mrasu/GravityR/cmd/flag"
	"github.com/mrasu/GravityR/cmd/util"
	"github.com/mrasu/GravityR/database/dmodel"
	"github.com/mrasu/GravityR/database/dservice"
	"github.com/mrasu/GravityR/database/hasura/hservice"
	"github.com/mrasu/GravityR/database/postgres/pmodel"
	"github.com/mrasu/GravityR/database/postgres/pservice"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/mrasu/GravityR/infra/hasura"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var HasuraCmd = &cobra.Command{
	Use:   "hasura",
	Short: "Suggest ways to increase Hasura's performance",
	Run: func(cmd *cobra.Command, args []string) {
		err := hasuraR.prepare()
		if err != nil {
			util.LogError(err)
			return
		}

		err = hasuraR.run()
		if err != nil {
			util.LogError(err)
		}
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

	parsedVariables map[string]interface{}
}

func (hr *hasuraRunner) prepare() error {
	v, err := hr.parseJSONToVariables(hr.jsonVariables)
	if err != nil {
		return err
	}

	hr.parsedVariables = v
	return nil
}

func (hr *hasuraRunner) run() error {
	cfg, err := hasura.NewConfigFromEnv()
	if err != nil {
		return err
	}
	cli := hasura.NewClient(cfg)

	return hr.suggest(flag.DbFlag.Output, cli)
}

func (hr *hasuraRunner) suggest(outputPath string, cli *hasura.Client) error {
	q := &hasura.ExplainRequestBody{
		Query: &hasura.Query{
			Query:     hr.query,
			Variables: hr.parsedVariables,
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
	aTree, err := pservice.NewExplainAnalyzeTreeBuilder().Build(r.Plan)
	if err != nil {
		return err
	}

	its, err := hservice.NewIndexSuggester(cli).Suggest(r.SQL, aTree)
	if err != nil {
		return err
	}

	examinationIdxTargets, err := parseIndexTargets(hr.indexTargets)
	if err != nil {
		return err
	}

	var er *dmodel.ExaminationResult
	if hr.runsExamination {
		er, err = hr.examine(cli, examinationIdxTargets, its)
		if err != nil {
			return err
		}
	}

	if outputPath != "" {
		err := hr.createHTML(outputPath, hr.parsedVariables, r.SQL, its, er, aTree)
		if err != nil {
			return err
		}

		util.LogResultOutputPath(outputPath)
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

func (hr *hasuraRunner) examine(cli *hasura.Client, varTargets, possibleTargets []*dmodel.IndexTarget) (*dmodel.ExaminationResult, error) {
	targets := varTargets
	if len(targets) == 0 {
		targets = lo.Filter(possibleTargets, func(it *dmodel.IndexTarget, _ int) bool { return it.IsSafe() })
	}

	log.Info().Msg("Start examination...")
	ie := hservice.NewIndexExaminer(cli, hr.query, hr.parsedVariables)
	er, err := dservice.NewIndexEfficiencyExaminer(ie).Run(targets)
	if err != nil {
		return nil, err
	}

	return er, nil
}

func (hr *hasuraRunner) createHTML(outputPath string, variables map[string]interface{}, sql string, idxTargets []*dmodel.IndexTarget, er *dmodel.ExaminationResult, aTree *pmodel.ExplainAnalyzeTree) error {
	vits := lo.Map(idxTargets, func(v *dmodel.IndexTarget, _ int) *viewmodel.VmIndexTarget { return v.ToViewModel() })

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
