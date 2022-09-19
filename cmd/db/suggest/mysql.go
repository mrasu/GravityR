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

type mysqlVarS struct {
	runsExamination bool
	indexTargets    []string
	query           string
}

var MySqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Suggest ways to increase MySQL's performance",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMySql()
	},
}

var mysqlVar = mysqlVarS{}

func init() {
	flg := MySqlCmd.Flags()
	flg.BoolVar(&mysqlVar.runsExamination, "with-examine", false, "Examine query by adding index")
	flg.StringArrayVarP(&mysqlVar.indexTargets, "index", "i", []string{}, "Specify index")

	flg.StringVarP(&mysqlVar.query, "query", "q", "", "[Required] Query to check")
	err := cobra.MarkFlagRequired(flg, "query")
	if err != nil {
		panic(err)
	}
}

func runMySql() error {
	/*	// MySQLに実在するテーブルのカラム一覧
		tables := map[string][]string{
			"users": {"user_id", "name", "email", "password", "created_at", "updated_at"},
			"todos": {"todo_id", "user_id", "title", "description", "status", "created_at", "updated_at"},
		}
		// クエリの各テーブルが参照しているカラムとテーブル
		columns := map[string][][]string{
			"count_tbl": {{"status(REFERENCE)", "count(GROUP)", "created_at(users)(CONDITION)", "user_id(users)(CONDITION)", "user_id(todos)(CONDITION)"}, {"users", "todos", "todos2"}},
			"haha":      {{"name(users)(REFERENCE)", "status(REFERENCE)", "user_id(users)(CONDITION)", "user_id(haha)(CONDITION)", "email(CONDITION)"}, {"users", "IAnalyzeData"}},
			"IAnalyzeData":         {{"status(REFERENCE)", "user_id(REFERENCE)"}, {"todos"}},
			// subqueryの場合はindexしても意味がないので、無視していい
			// subqueryの後は集計結果を使って集計するから、indexはFROMが直接テーブルを参照しているテーブルのみに関係し、subquery対象のwhereは意味がない
			// 		(where して絞り込めばoutputする行数が減るから効果はあるが、indexをつけるのとは別の改善方法 (クエリの書き換え)なので今は対象外
			"root":      {{"status(CONDITION)", "count(CONDITION)"}, {"count_tbl"}},
		}

		// 各テーブルが使用しているカラムの元のテーブル名
		IAnalyzeData := map[string]map[string][]string{
			"count_tbl": {
				"users":  {"user_id(CONDITION)", "created_at(CONDITION)"},
				"todos":  {"user_id(CONDITION)", "status(REFERENCE)"},
				"todos2": {"user_id(CONDITION)"},
			},
			"haha": {
				"users":      {"user_id(CONDITION)", "name(REFERENCE)", "email(CONDITION)"},
				"todos(sub)": {"status(REFERENCE)", "user_id(CONDITION)"},
			},
			"IAnalyzeData": {
				"todos": {"status(REFERENCE)", "user_id(REFERENCE)"},
			},
			"root": {
				"todos": {"status(CONDITION)"},
				"GROUP": {"count(CONDITION)"},
			},
		}
		// インデックスをかけると意味がありそうなカラムの組み合わせ(順不同)
		idxCandidates := map[string][][]string{
			"users": {
				{"user_id(CONDITION)", "created_at(CONDITION)"},
				{"user_id(CONDITION)", "name(REFERENCE)", "email(CONDITION)"},
			},
			"todos": {
				{"user_id(CONDITION)", "status(REFERENCE)"},
				{"status(REFERENCE)", "user_id(CONDITION)"},
				{"status(CONDITION)"},
			},
			"todos2": {
				{"user_id(CONDITION)"},
			},
		}
		// 順序も考慮されたインデックスの候補
		// これだけ候補があるが、bottomなselectは `count_tbl`と`IAnalyzeData`であるから、コスト推計はその2つだけサポートする (未来は知らん)
		// 他の候補は表示だけ。
		// 表示する時に、「このインデックスをつければこの部分が改善される」みたいなのを出せるとわかりやすい (直感に反するインデックス範囲があるので)
		// REFERENCEより、CONDITIONの方が価値があるので、CONDITION用のインデックスの方が可能性が高い
		idxFullOrderedCandidates := map[string][][]string{
			"users": {
				{"created_at"}, // 実質、{created_at, user_id}だが、user_idはprimaryだから最後につける必要はない。
				{"user_id", "created_at"},
				{"email"},            // 実質、{email, user_id}だが、user_idはprimaryだから最後につける必要はない。
				{"user_id", "email"}, // スキャンで絞れているならnameをつける必要はない (nameはインデックススキャンで終わるようにするためのもの)
				{"user_id", "email", "name"},
				{"user_id", "email", "name"},
			},
			"todos": {
				{"user_id", "status"},
				{"status", "user_id"},
				{"status"}, // {status, user_id}と似ているが、statusで絞れるならuser_idのインデックスを使う必要はない
			},
			"todos2": {
				{"user_id"},
			},
		}
		costs := map[string]string{
			"users": "111",
			"todos": "222",
			"count_tbl": "5"
			"count_tbl": "333"
		}

		// 各テーブルが使用しているカラムの元のテーブル名(no-subquery)
		a2 := map[string]map[string][]string{
			"count_tbl": {
				"users":  {"user_id(CONDITION)", "created_at(CONDITION)"},
				"todos":  {"user_id(CONDITION)", "status(REFERENCE)"},
				"todos2": {"user_id(CONDITION)"},
			},
			"haha": {
				"users":      {"user_id(CONDITION)", "name(REFERENCE)", "email(CONDITION)"},
			},
			"IAnalyzeData": {
				"todos": {"status(REFERENCE)", "user_id(REFERENCE)"},
			}
		}
		// インデックスをかけると意味がありそうなカラムの組み合わせ(順不同)(no-subquery)
		idxCandidates2 := map[string][][]string{
			"users": {
				"count_tbl": {"user_id(CONDITION)", "created_at(CONDITION)"},
				"haha": {"user_id(CONDITION)", "name(REFERENCE)", "email(CONDITION)"},
			},
			"todos": {
				"count_tbl": {"user_id(CONDITION)", "status(REFERENCE)"},
			},
			"todos2": {
				"count_tbl": {"user_id(CONDITION)"},
			},
		}
		// 順序も考慮されたインデックスの候補
		idxFullOrderedCandidates2 := map[string][][]string{
			"users": {
				"count_tbl": {
					{"created_at"}, // 実質、{created_at, user_id}だが、user_idはprimaryだから最後につける必要はない。
					{"user_id", "created_at"},
				},
				"haha": {
					{"email"},            // 実質、{email, user_id}だが、user_idはprimaryだから最後につける必要はない。
					{"user_id", "email"}, // スキャンで絞れているならnameをつける必要はない (nameはインデックススキャンで終わるようにするためのもの)
					{"user_id", "email", "name"},
					{"email", "user_id", "name"},
				}
			},
			"todos": {
				"count_tbl": {
					{"user_id"},
					{"user_id", "status"},
				}
			},
			"todos2": {
				"count_tbl": {{"user_id"},}
			},
		}
		costs2 := map[string]string{
			"users": "111",
			"todos": "222",
			"todos2": "3",
		}
		// カーディナリティ計算 (CONDITIONだけが対象。REFERENCEはindexの絞り込み効果に関係がない)
		cardsQuery2 := {
			"users": `SELECT
				COUNT(DISTINCT created_at) / COUNT(*),
				COUNT(DISTINCT user_id, created_at) / COUNT(*),
				COUNT(DISTINCT email LIKE 'test%') / COUNT(*),
				COUNT(DISTINCT user_id, email) / COUNT(*),
				COUNT(DISTINCT email, user_id) / COUNT(*),
			"todos": `
				COUNT(DISTINCT user_id) / COUNT(*),
			`
			"todos2": `
				COUNT(DISTINCT user_id) / COUNT(*),
			`
		}
		costs2が大きく、cardQuery2の結果が大きいものは効果が高い。
		効果が高いもの+REFERENCEの組み合わせが上位の「やると意味あるかも」なインデックス

		=> この後に実際にインデックスを貼って効果検証するフェーズが来るが、それは別ファイル
	*/

	/*
			query := `
		with count_tbl as
		( select
		  status,
		  count(status) as count
		  from users
		  inner join todos on users.user_id = todos.user_id
		  where users.created_at > now() - interval 2  month
		  group by status
		),
		haha as (
		 select users.name, status from users inner join (select status, user_id from todos) as IAnalyzeData on users.user_id = IAnalyzeData.user_id
		 where email like 'test%'
		)
		select
		  case when status = 1 then 'Todo'
		       when status = 2 then 'Doing'
		       when status = 3 then 'Done' end as status,
		  count / ( select sum(count) from count_tbl ) * 100 as percent
		from  count_tbl
		where status > 0 and count > 0
		order by status desc;
		`
	*/

	/*
			query := `
		SELECT
			name,
			t.description
		FROM
			users
			INNER JOIN todos AS t ON users.id = t.user_id
		WHERE
			t.user_id < 1000
		`
	*/
	fmt.Printf("\nQuery: %s\n\n", mysqlVar.query)

	cfg, err := iMysql.NewConfigFromEnv()
	if err != nil {
		return err
	}

	db, err := iMysql.OpenMySQLDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	/*
		idxTargets = []*model.IndexTargetTable{
			{
				TableName: "todos",
				IndexFields: []*model.IndexField{
					{Name: "user_id"},
				},
			},
			{
				TableName: "todos",
				IndexFields: []*model.IndexField{
					{Name: "status"},
				},
			},
			{
				TableName: "todos",
				IndexFields: []*model.IndexField{
					{Name: "user_id"},
					{Name: "status"},
					{Name: "description"},
				},
			},
		}
	*/

	examinationIdxTargets, err := parseIndexTargets(mysqlVar.indexTargets)
	if err != nil {
		return err
	}

	explainLine, err := db.Explain(mysqlVar.query)
	if err != nil {
		return err
	}

	aTree, err := collector.CollectExplainAnalyzeTree(explainLine)
	if err != nil {
		return err
	}

	its, errs := mysql.SuggestIndex(db, cfg.GetDBName(), mysqlVar.query, aTree)
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
	if mysqlVar.runsExamination {
		fmt.Printf("\n======going to examine-------\n")
		ie := mysql.NewIndexExaminer(db, mysqlVar.query)
		er, err = database.NewIndexEfficiencyExaminer(ie).Run(examinationIdxTargets)
		if err != nil {
			return err
		}
	}

	outputPath := flag.DbFlag.Output
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
			mysqlVar.query,
			aTree.ToViewModel(),
			vits,
			[]*viewmodel.VmExaminationCommandOption{
				viewmodel.CreateOutputExaminationOption(!mysqlVar.runsExamination, outputPath),
				{IsShort: true, Name: "q", Value: mysqlVar.query},
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
