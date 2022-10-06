package service

import "github.com/mrasu/GravityR/database"

type ExplainNode interface {
	TableName() string
	EstimatedTotalTime() float64

	GetChildren() []ExplainNode
}

type singleTableTree struct {
	tableName string
	nodes     []ExplainNode
}

type singleTableExplainResultsBuilder struct{}

func BuildSingleTableExplainResults(root ExplainNode) []*database.SingleTableExplainResult {
	var trees []*singleTableTree
	sb := singleTableExplainResultsBuilder{}
	sb.toSingleTableTreeRecursive(&trees, root, nil)

	var res []*database.SingleTableExplainResult
	for _, tree := range trees {
		res = append(res, &database.SingleTableExplainResult{
			TableName:          tree.tableName,
			EstimatedTotalTime: tree.nodes[0].EstimatedTotalTime(),
		})
	}
	return res
}

func (sb *singleTableExplainResultsBuilder) toSingleTableTreeRecursive(trees *[]*singleTableTree, node ExplainNode, currentTree *singleTableTree) bool {
	treeUsed := false
	isRoot := false
	if currentTree == nil {
		if node.TableName() != "" {
			currentTree = &singleTableTree{
				tableName: node.TableName(),
				nodes:     []ExplainNode{node},
			}
			isRoot = true
		}
	} else {
		if node.TableName() == currentTree.tableName {
			currentTree.nodes = append(currentTree.nodes, node)
			treeUsed = true
		} else {
			*trees = append(*trees, currentTree)

			if node.TableName() != "" {
				currentTree = &singleTableTree{
					tableName: node.TableName(),
					nodes:     []ExplainNode{node},
				}
				isRoot = true
			} else {
				currentTree = nil
			}
		}
	}

	if len(node.GetChildren()) == 0 && isRoot {
		*trees = append(*trees, currentTree)
		return false
	}

	for _, c := range node.GetChildren() {
		used := sb.toSingleTableTreeRecursive(trees, c, currentTree)
		if used {
			if isRoot {
				*trees = append(*trees, currentTree)
			}
			currentTree = nil
			treeUsed = true
		}
	}

	return treeUsed
}
