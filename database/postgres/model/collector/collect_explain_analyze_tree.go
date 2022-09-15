package collector

import (
	"github.com/mrasu/GravityR/database/postgres/model"
	"github.com/mrasu/GravityR/lib"
	"regexp"
	"strings"
)

type explainAnalyzeResultCollector struct{}

func CollectExplainAnalyzeTree(explainLines []string) (*model.ExplainAnalyzeTree, error) {
	c := explainAnalyzeResultCollector{}
	root, summaryText, err := c.collect(explainLines)
	if err != nil {
		return nil, err
	}

	return &model.ExplainAnalyzeTree{
		Root:        root,
		SummaryText: summaryText,
	}, nil
}

func (earc *explainAnalyzeResultCollector) collect(explainLines []string) (*model.ExplainAnalyzeTreeNode, string, error) {
	lineNodes, summaryText := earc.groupToPlanLineNodes(explainLines)

	root := &model.ExplainAnalyzeTreeNode{AnalyzeResultNode: &model.ExplainAnalyzeResultNode{}, SpaceSize: -1}
	nodeStack := lib.NewStack[model.ExplainAnalyzeTreeNode]()
	nodeStack.Push(root)
	for _, ln := range lineNodes {
		spaceSize, lNode, err := model.ParseExplainAnalyzeResultLineNode(ln)
		if err != nil {
			return nil, "", err
		}
		n := &model.ExplainAnalyzeTreeNode{
			AnalyzeResultNode: lNode,
			Children:          nil,
			SpaceSize:         spaceSize,
		}
		if nodeStack.Top().SpaceSize < spaceSize {
			currentNode := nodeStack.Top()
			currentNode.Children = append(currentNode.Children, n)
		} else {
			nodeStack.Pop()
			for nodeStack.Top().SpaceSize >= spaceSize {
				nodeStack.Pop()
			}
			currentNode := nodeStack.Top()
			currentNode.Children = append(currentNode.Children, n)
		}

		nodeStack.Push(n)
	}
	return root, summaryText, nil
}

var arrowRowReg = regexp.MustCompile(`\s*-> `)
var infoRowReg = regexp.MustCompile(`\s*(\w+ ?)+: `)
var summaryReg = regexp.MustCompile(`^\w[\w ]+:`)

const (
	analyzeLineInfo int = iota
	analyzeLineArrow
	analyzeLinePlanName
)

func (earc *explainAnalyzeResultCollector) groupToPlanLineNodes(explainLines []string) ([][]string, string) {
	var lineNodes [][]string
	var nodeLines []string
	prevLineType := analyzeLinePlanName
	summaryStartLine := 0
	for i, line := range explainLines {
		if summaryReg.MatchString(line) {
			if len(nodeLines) > 0 {
				lineNodes = append(lineNodes, nodeLines)
			}
			summaryStartLine = i
			break
		}
		if i == 0 {
			nodeLines = append(nodeLines, line)
			prevLineType = analyzeLineArrow
			continue
		}
		if prevLineType == analyzeLineInfo {
			if infoRowReg.MatchString(line) {
				nodeLines = append(nodeLines, line)
				continue
			} else {
				if len(nodeLines) > 0 {
					lineNodes = append(lineNodes, nodeLines)
					nodeLines = nil
				}
			}
		}
		if infoRowReg.MatchString(line) {
			nodeLines = append(nodeLines, line)
			prevLineType = analyzeLineInfo
		} else if arrowRowReg.MatchString(line) {
			if prevLineType == analyzeLineArrow {
				if len(nodeLines) > 0 {
					lineNodes = append(lineNodes, nodeLines)
					nodeLines = nil
				}
			}
			nodeLines = append(nodeLines, line)
			prevLineType = analyzeLineArrow
		} else {
			if len(nodeLines) > 0 {
				lineNodes = append(lineNodes, nodeLines)
				nodeLines = nil
			}
			nodeLines = append(nodeLines, line)
			prevLineType = analyzeLinePlanName
		}
	}

	sLine := ""
	if summaryStartLine > 0 {
		sLine = strings.Join(explainLines[summaryStartLine:], "\n")
	}

	return lineNodes, sLine
}
