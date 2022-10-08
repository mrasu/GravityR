package parser

import (
	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/pkg/errors"
)

func Parse(query string) (*parser.Statement, error) {
	stmts, err := parser.Parse(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse sql")
	}

	if len(stmts) > 1 {
		return nil, errors.New("not supporting query having multiple statements")
	}
	return &stmts[0], nil
}
