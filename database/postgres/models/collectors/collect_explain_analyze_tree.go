package collectors

import (
	"github.com/jmoiron/sqlx"
	"github.com/mrasu/GravityR/database/postgres/models"
	"github.com/mrasu/GravityR/lib"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

type explainAnalyzeResultCollector struct{}

func CollectExplainAnalyzeTree(db *sqlx.DB, query string) (*models.ExplainAnalyzeTree, error) {
	c := explainAnalyzeResultCollector{}
	root, planningText, err := c.collect(db, query)
	if err != nil {
		return nil, err
	}

	return &models.ExplainAnalyzeTree{
		Root:         root,
		PlanningText: planningText,
	}, nil
}

func (earc *explainAnalyzeResultCollector) collect(db *sqlx.DB, query string) (*models.ExplainAnalyzeTreeNode, string, error) {
	treeLines, err := earc.runExplainAnalyzeResult(db, query)
	if err != nil {
		return nil, "", err
	}
	lineNodes, planningText := earc.groupToPlanLineNodes(treeLines)

	root := &models.ExplainAnalyzeTreeNode{AnalyzeResultNode: &models.ExplainAnalyzeResultNode{}, SpaceSize: -1}
	nodeStack := lib.NewStack[models.ExplainAnalyzeTreeNode]()
	nodeStack.Push(root)
	for _, ln := range lineNodes {
		spaceSize, lNode, err := models.ParseExplainAnalyzeResultLineNode(ln)
		if err != nil {
			return nil, "", err
		}
		n := &models.ExplainAnalyzeTreeNode{
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
	return root, planningText, nil
}

func (earc *explainAnalyzeResultCollector) runExplainAnalyzeResult(db *sqlx.DB, query string) ([]string, error) {
	rows, err := db.Query("EXPLAIN (ANALYZE, BUFFERS) " + query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to select")
	}
	defer rows.Close()

	var res []string
	for rows.Next() {
		var txt string
		if err := rows.Scan(&txt); err != nil {
			return nil, errors.Wrap(err, "failed to Scan")
		}
		res = append(res, txt)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to execute EXPLAIN ANALYZE")
	}

	return res, nil
}

var arrowRowReg = regexp.MustCompile(`\s*-> `)
var infoRowReg = regexp.MustCompile(`\s*(\w+ ?)+: `)

const (
	analyzeLineInfo int = iota
	analyzeLineArrow
	analyzeLineOther
)

func (earc *explainAnalyzeResultCollector) groupToPlanLineNodes(explainLines []string) ([][]string, string) {
	var lineNodes [][]string
	var nodeLines []string
	prevLineType := analyzeLineOther
	pStartLine := 0
	for i, line := range explainLines {
		// `^Planning` is a marker which indicates following lines are summary for the query
		if strings.HasPrefix(line, "Planning") {
			if len(nodeLines) > 0 {
				lineNodes = append(lineNodes, nodeLines)
			}
			pStartLine = i
			break
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
			nodeLines = append(nodeLines, line)
			prevLineType = analyzeLineOther
		}
	}

	pLine := ""
	if pStartLine > 0 {
		pLine = strings.Join(explainLines[pStartLine:], "\n")
	}

	return lineNodes, pLine
}
