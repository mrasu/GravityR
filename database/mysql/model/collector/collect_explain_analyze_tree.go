package collector

import (
	"github.com/mrasu/GravityR/database/mysql/model"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"strings"
)

type explainAnalyzeResultCollector struct{}

func CollectExplainAnalyzeTree(explainLine string) (*model.ExplainAnalyzeTree, error) {
	c := explainAnalyzeResultCollector{}
	root, err := c.collect(explainLine)
	if err != nil {
		return nil, err
	}

	return &model.ExplainAnalyzeTree{Root: root}, nil
}

func (earc *explainAnalyzeResultCollector) collect(explainLine string) (*model.ExplainAnalyzeTreeNode, error) {
	nodeStack := lib.NewStack[model.ExplainAnalyzeTreeNode]()

	root := &model.ExplainAnalyzeTreeNode{AnalyzeResultLine: &model.ExplainAnalyzeResultLine{}}
	nodeStack.Push(root)
	treeLines := strings.Split(explainLine, "\n")
	for _, line := range treeLines {
		if line == "" {
			continue
		}

		nest, l, err := model.ParseExplainAnalyzeResultLine(line)
		if err != nil {
			return nil, err
		}
		n := &model.ExplainAnalyzeTreeNode{
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
