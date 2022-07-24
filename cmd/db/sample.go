package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mrasu/GravityR/database/mysql"
	"github.com/mrasu/GravityR/database/mysql/models"
	"github.com/mrasu/GravityR/database/mysql/models/collectors"
	"github.com/mrasu/GravityR/html"
	"github.com/mrasu/GravityR/html/viewmodel"
	"github.com/mrasu/GravityR/lib"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/spf13/cobra"
)

// SampleCmd represents the sample command
var SampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is IAnalyzeData CLI library for Go that empowers applications.
This application is IAnalyzeData tool to generate the needed files
to quickly create IAnalyzeData Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSample()
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// suggestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// suggestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runSample() error {
	fmt.Println("Running sample...")

	cfg, err := mysql.NewConfigFromEnv()
	if err != nil {
		return err
	}

	db, err := mysql.OpenMySQLDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	query := "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"

	rootNode, err := parse(query)
	if err != nil {
		panic(err)
	}
	tNames, errs := collectors.CollectTableNames(rootNode)
	if len(errs) > 0 {
		return errs[0]
	}

	tables, err := collectors.CollectTableSchemas(db, "gravityr", tNames)
	if err != nil {
		return err
	}

	scopes, errs := collectors.CollectStmtScopes(rootNode)
	if len(errs) > 0 {
		return errs[0]
	}

	idxCandidates, err := collectors.CollectIndexTargets(tables, scopes)
	if err != nil {
		return err
	}
	fmt.Println(lib.JoinF(idxCandidates, "\n", func(f *models.IndexTargetTable) string { return f.String() }))
	fmt.Println()
	var vits []*viewmodel.VmIndexTarget
	for _, it := range idxCandidates {
		vits = append(vits, it.ToIndexTarget().ToViewModel())
	}

	aTree, err := collectors.CollectExplainAnalyzeTree(db, query)
	if err != nil {
		return err
	}

	bo := &html.BuildOption{
		Query:        query,
		AnalyzeNodes: aTree.ToViewModel(),
		IndexTargets: vits,
		CommandOptions: []*viewmodel.VmExaminationCommandOption{
			viewmodel.CreateOutputExaminationOption(true, "output.html"),
			{IsShort: true, Name: "q", Value: query},
		},
		ExaminationResults: nil,
	}

	err = html.CreateHtml("output.html", bo)
	if err != nil {
		return err
	}

	return nil
}

func parse(sql string) (ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return stmtNodes[0], nil
}
