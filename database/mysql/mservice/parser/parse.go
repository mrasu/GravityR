package parser

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
)

func Parse(query string) (ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(query, "", "")
	if err != nil {
		return nil, err
	}

	return stmtNodes[0], nil
}
