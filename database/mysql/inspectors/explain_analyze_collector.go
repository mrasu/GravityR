package inspectors

import (
	"github.com/jmoiron/sqlx"
	"github.com/mrasu/GravityR/database/mysql/models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"strings"
)

type explainAnalyzeResultCollector struct{}

func CollectExplainAnalyzeTree(db *sqlx.DB, query string) (*models.ExplainAnalyzeTree, error) {
	nodeStack := lib.NewStack[models.ExplainAnalyzeTreeNode]()

	tree, err := getExplainAnalyzeResult(db, query)
	if err != nil {
		return nil, err
	}

	root := &models.ExplainAnalyzeTreeNode{AnalyzeResultLine: &models.ExplainAnalyzeResultLine{}}
	nodeStack.Push(root)
	treeLines := strings.Split(tree, "\n")
	for _, line := range treeLines {
		if line == "" {
			continue
		}

		nest, l, err := models.ParseExplainAnalyzeResultLine(line)
		if err != nil {
			return nil, err
		}
		n := &models.ExplainAnalyzeTreeNode{
			AnalyzeResultLine: l,
			Children:          nil,
		}
		if nodeStack.Size()-1 == nest {
			currentNode := nodeStack.Top()
			currentNode.Children = append(currentNode.Children, n)
		} else if nodeStack.Size()-1 > nest {
			nodeStack.Pop()
			for nodeStack.Size()-1 > nest {
				nodeStack.Pop()
			}
			currentNode := nodeStack.Top()
			currentNode.Children = append(currentNode.Children, n)
		} else {
			return nil, errors.New("invalid result from EXPLAIN ANALYZE")
		}

		nodeStack.Push(n)
	}

	tmp := &models.ExplainAnalyzeTree{Root: root}
	return tmp, nil
}

func getExplainAnalyzeResult(db *sqlx.DB, query string) (string, error) {
	rows, err := db.Query("EXPLAIN ANALYZE FORMAT=TREE " + query)
	if err != nil {
		return "", errors.Wrap(err, "failed to select")
	}
	defer rows.Close()

	var res string
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return "", errors.Wrap(err, "failed to execute EXPLAIN ANALYZE")
		} else {
			return "", errors.New("no result for EXPLAIN ANALYZE")
		}
	}

	if err := rows.Scan(&res); err != nil {
		return "", errors.Wrap(err, "failed to Scan")
	}

	return res, nil
}
