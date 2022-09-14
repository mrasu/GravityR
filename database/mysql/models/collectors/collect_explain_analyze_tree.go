package collectors

import (
	"github.com/mrasu/GravityR/database/mysql/models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"strings"
)

type explainAnalyzeResultCollector struct{}

func CollectExplainAnalyzeTree(explainLine string) (*models.ExplainAnalyzeTree, error) {
	c := explainAnalyzeResultCollector{}
	root, err := c.collect(explainLine)
	if err != nil {
		return nil, err
	}

	return &models.ExplainAnalyzeTree{Root: root}, nil
}

func (earc *explainAnalyzeResultCollector) collect(explainLine string) (*models.ExplainAnalyzeTreeNode, error) {
	nodeStack := lib.NewStack[models.ExplainAnalyzeTreeNode]()

	root := &models.ExplainAnalyzeTreeNode{AnalyzeResultLine: &models.ExplainAnalyzeResultLine{}}
	nodeStack.Push(root)
	treeLines := strings.Split(explainLine, "\n")
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

	return root, nil
}
