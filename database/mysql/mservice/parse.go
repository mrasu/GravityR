package mservice

import (
	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
)

func Parse(p *parser.Parser, query string) (ast.StmtNode, error) {
	stmtNodes, _, err := p.Parse(query, "", "")
	if err != nil {
		return nil, err
	}

	return stmtNodes[0], nil
}
