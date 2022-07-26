package hasura

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

var PostgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Suggest ways to increase performance of Hasura with PostgreSQL",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := postgresR.prepare()
		if err != nil {
			return err
		}

		err = postgresR.run()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	flg := PostgresCmd.Flags()
	flg.BoolVar(&postgresR.runsExamination, "with-examine", false, "Examine query by adding index")
	flg.StringArrayVarP(&postgresR.indexTargets, "index", "i", []string{}, "Specify index")

	flg.StringVarP(&postgresR.query, "query", "q", "", "[Required] Query of GraphQL to check")
	flg.StringVar(&postgresR.jsonVariables, "json-variables", "{}", "JSON of variables for the query")
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
	jsonVariables   string

	parsedVariables map[string]interface{}
}

func (pr *postgresRunner) prepare() error {
	v, err := pr.parseJSONToVariables(pr.jsonVariables)
	if err != nil {
		return err
	}

	pr.parsedVariables = v
	return nil
}

func (pr *postgresRunner) run() error {
	cfg, err := hasura.NewConfigFromEnv()
	if err != nil {
		return err
	}
	cli := hasura.NewClient(cfg)

	return pr.suggest(flag.DbFlag.Output, cli)
}

func (pr *postgresRunner) suggest(outputPath string, cli *hasura.Client) error {
	q := &hasura.ExplainRequestBody{
		Query: &hasura.Query{
			Query:     pr.query,
			Variables: pr.parsedVariables,
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

	its, err := hservice.NewIndexSuggester(cli).Suggest(r.SQL)
	if err != nil {
		return err
	}
	util.LogNewIndexTargets(its)

	examinationIdxTargets, err := dmodel.NewIndexTargetsFromTexts(pr.indexTargets)
	if err != nil {
		return err
	}

	var er *dmodel.ExaminationResult
	if pr.runsExamination {
		er, err = pr.examine(cli, examinationIdxTargets, its)
		if err != nil {
			return err
		}
	}

	if outputPath != "" {
		err := pr.createHTML(outputPath, pr.parsedVariables, r.SQL, its, er, aTree)
		if err != nil {
			return err
		}

		util.LogResultOutputPath(outputPath)
	}
	return nil
}

func (pr *postgresRunner) parseJSONToVariables(jsonStr string) (map[string]interface{}, error) {
	var variables map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &variables); err != nil {
		return nil, errors.Wrap(err, "failed to parse variables")
	}

	return variables, nil
}

func (pr *postgresRunner) examine(cli *hasura.Client, varTargets, possibleTargets []*dmodel.IndexTarget) (*dmodel.ExaminationResult, error) {
	targets := varTargets
	if len(targets) == 0 {
		targets = lo.Filter(possibleTargets, func(it *dmodel.IndexTarget, _ int) bool { return it.IsSafe() })
	}

	log.Info().Msg("Start examination...")
	ie := hservice.NewIndexExaminer(cli, pr.query, pr.parsedVariables)
	er, err := dservice.NewIndexEfficiencyExaminer(ie).Run(targets)
	if err != nil {
		return nil, err
	}

	return er, nil
}

func (pr *postgresRunner) createHTML(outputPath string, variables map[string]interface{}, sql string, idxTargets []*dmodel.IndexTarget, er *dmodel.ExaminationResult, aTree *pmodel.ExplainAnalyzeTree) error {
	vits := lo.Map(idxTargets, func(v *dmodel.IndexTarget, _ int) *viewmodel.VmIndexTarget { return v.ToViewModel() })

	var ver *viewmodel.VmExaminationResult
	if er != nil {
		ver = er.ToViewModel()
	}

	bo := html.NewSuggestHasuraDataBuildOption(
		pr.query,
		variables,
		sql,
		aTree.ToViewModel(),
		aTree.SummaryText,
		vits,
		[]*viewmodel.VmExaminationCommandOption{
			viewmodel.CreateOutputExaminationOption(!pr.runsExamination, outputPath),
			{IsShort: true, Name: "q", Value: pr.query},
			{IsShort: false, Name: "json-variables", Value: pr.jsonVariables},
		},
		ver,
	)

	err := html.CreateHtml(html.TypeMain, outputPath, bo)
	if err != nil {
		return err
	}

	return nil
}
