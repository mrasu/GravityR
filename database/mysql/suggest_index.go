package mysql

import (
	"fmt"
	"github.com/mrasu/GravityR/database/common_model"
	"github.com/mrasu/GravityR/database/common_model/builder"
	"github.com/mrasu/GravityR/database/mysql/model"
	"github.com/mrasu/GravityR/database/mysql/model/collector"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/mrasu/GravityR/lib"
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
)

/*
	query := `with count_tbl as
	( select
	  status,
	  count(status) as count
	  from users
	  inner join todos on users.user_id = todos.user_id
	  where users.created_at > now() - interval 2  month
	  group by status
	),
	haha as (
	 select users.name, status from users inner join (select status, user_id from todos) as a on users.user_id = a.user_id
	 where email like 'test%'
	)
	select
	  case when status = 1 then 'Tod'
		   when status = 2 then 'Doing'
		   when status = 3 then 'Done' end as status,
	  count / ( select sum(count) from count_tbl ) * 100 as percent
	from  count_tbl
	where status > 0 and count > 0
	order by status desc;`

	// MySQLに実在するテーブルのカラム一覧
	tables := map[string][]string{
		"users": {"user_id", "name", "email", "password", "created_at", "updated_at"},
		"todos": {"todo_id", "user_id", "title", "description", "status", "created_at", "updated_at"},
	}
	// クエリの各テーブルが参照しているカラムとテーブル
	columns := map[string][][]string{
		"count_tbl": {{"status(REFERENCE)", "count(GROUP)", "created_at(users)(CONDITION)", "user_id(users)(CONDITION)", "user_id(todos)(CONDITION)"}, {"users", "todos", "todos2"}},
		"haha":      {{"name(users)(REFERENCE)", "status(REFERENCE)", "user_id(users)(CONDITION)", "user_id(haha)(CONDITION)", "email(CONDITION)"}, {"users", "a"}},
		"a":         {{"status(REFERENCE)", "user_id(REFERENCE)"}, {"todos"}},
		// subqueryの場合はindexしても意味がないので、無視していい
		// subqueryの後は集計結果を使って集計するから、indexはFROMが直接テーブルを参照しているテーブルのみに関係し、subquery対象のwhereは意味がない
		// 		(where して絞り込めばoutputする行数が減るから効果はあるが、indexをつけるのとは別の改善方法 (クエリの書き換え)なので今は対象外
		"root":      {{"status(CONDITION)", "count(CONDITION)"}, {"count_tbl"}},
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
		"a": {
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
	// これだけ候補があるが、bottomなselectは `count_tbl`と`a`であるから、コスト推計はその2つだけサポートする (未来は知らん)
	// 他の候補は表示だけ。
	// 表示する時に、「このインデックスをつければこの部分が改善される」みたいなのを出せるとわかりやすい (直感に反するインデックス範囲があるので)
	// REFERENCEより、CONDITIONの方が価値があるので、CONDITION用のインデックスの方が可能性が高い
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

// TODO: 既存のインデックスと被るものは除外する
func SuggestIndex(db *mysql.DB, database, query string, aTree *model.ExplainAnalyzeTree) ([]*common_model.IndexTargetTable, []error) {
	rootNode, err := parse(query)
	if err != nil {
		return nil, []error{err}
	}
	tNames, errs := collector.CollectTableNames(rootNode)
	if len(errs) > 0 {
		return nil, errs
	}

	tables, err := collector.CollectTableSchemas(db, database, tNames)
	if err != nil {
		return nil, []error{err}
	}

	scopes, errs := collector.CollectStmtScopes(rootNode)
	if len(errs) > 0 {
		return nil, errs
	}

	idxCandidates, err := builder.BuildIndexTargets(tables, scopes)
	if err != nil {
		return nil, []error{err}
	}
	fmt.Println(lib.Join(idxCandidates, "\n", func(f *common_model.IndexTargetTable) string { return f.String() }))
	fmt.Println()

	tableResults := aTree.ToSingleTableResults()
	fmt.Println(lib.Join(tableResults, "\n", func(st *common_model.SingleTableExplainResult) string { return st.String() }))

	return builder.BuildExplainedIndexTargets(idxCandidates, scopes, tableResults)
}

func parse(sql string) (ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return stmtNodes[0], nil
}
