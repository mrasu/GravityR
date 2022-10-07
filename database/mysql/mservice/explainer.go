package mservice

import (
	"github.com/mrasu/GravityR/database/mysql/mmodel"
	"github.com/mrasu/GravityR/infra/mysql"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"strings"
)

type Explainer struct {
	db *mysql.DB
}

func NewExplainer(db *mysql.DB) *Explainer {
	return &Explainer{db: db}
}

func (e *Explainer) ExplainWithAnalyze(query string) (*mmodel.ExplainAnalyzeTree, error) {
	explainLine, err := e.db.ExplainWithAnalyze(query)
	if err != nil {
		return nil, err
	}

	root, err := e.buildExplainNode(explainLine)
	if err != nil {
		return nil, err
	}

	return &mmodel.ExplainAnalyzeTree{Root: root}, nil
}

func (e *Explainer) buildExplainNode(explainLine string) (*mmodel.ExplainAnalyzeTreeNode, error) {
	nodeStack := lib.NewStack[mmodel.ExplainAnalyzeTreeNode]()

	root := &mmodel.ExplainAnalyzeTreeNode{AnalyzeResultLine: &mmodel.ExplainAnalyzeResultLine{}}
	nodeStack.Push(root)
	treeLines := strings.Split(explainLine, "\n")
	for _, line := range treeLines {
		if line == "" {
			continue
		}

		nest, l, err := mmodel.ParseExplainAnalyzeResultLine(line)
		if err != nil {
			return nil, err
		}
		n := &mmodel.ExplainAnalyzeTreeNode{
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

	return root, nil
}
