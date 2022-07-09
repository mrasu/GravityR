package models

import "github.com/mrasu/GravityR/html/viewmodel"

type ExplainAnalyzeTree struct {
	Root *ExplainAnalyzeTreeNode
}

type singleTableTree struct {
	tableName string
	nodes     []*ExplainAnalyzeTreeNode
}

func (eat *ExplainAnalyzeTree) ToSingleTableTrees() []*SingleTableExplainAnalyzeTree {
	var trees []*singleTableTree
	eat.toSingleTableTreeRecursive(&trees, eat.Root, nil)

	var res []*SingleTableExplainAnalyzeTree
	for _, tree := range trees {
		res = append(res, &SingleTableExplainAnalyzeTree{
			TableName:          tree.tableName,
			EstimatedTotalTime: tree.nodes[0].AnalyzeResultLine.EstimatedTotalTime(),
		})
	}
	return res
}

func (eat *ExplainAnalyzeTree) toSingleTableTreeRecursive(trees *[]*singleTableTree, node *ExplainAnalyzeTreeNode, currentTree *singleTableTree) bool {
	treeUsed := false
	isRoot := false
	if currentTree == nil {
		if node.AnalyzeResultLine.TableName != "" {
			currentTree = &singleTableTree{
				tableName: node.AnalyzeResultLine.TableName,
				nodes:     []*ExplainAnalyzeTreeNode{node},
			}
			isRoot = true
		}
	} else {
		if node.AnalyzeResultLine.TableName == currentTree.tableName {
			currentTree.nodes = append(currentTree.nodes, node)
			treeUsed = true
		} else {
			*trees = append(*trees, currentTree)

			if node.AnalyzeResultLine.TableName != "" {
				currentTree = &singleTableTree{
					tableName: node.AnalyzeResultLine.TableName,
					nodes:     []*ExplainAnalyzeTreeNode{node},
				}
				isRoot = true
			} else {
				currentTree = nil
			}
		}
	}

	if len(node.Children) == 0 && isRoot {
		*trees = append(*trees, currentTree)
		return false
	}

	for _, c := range node.Children {
		used := eat.toSingleTableTreeRecursive(trees, c, currentTree)
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

func (eat *ExplainAnalyzeTree) ToViewModel() []*viewmodel.VmExplainAnalyzeNode {
	var res []*viewmodel.VmExplainAnalyzeNode
	for _, n := range eat.Root.Children {
		res = append(res, n.ToViewModel())
	}
	return res
}
