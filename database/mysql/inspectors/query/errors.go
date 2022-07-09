package query

import (
	"fmt"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pkg/errors"
	"reflect"
)

type UnexpectedChildError struct {
	parent ast.Node
	child  ast.Node
}

func NewUnexpectedChildError(p, c ast.Node) error {
	return errors.Wrap(&UnexpectedChildError{p, c}, "")
}

func (uc *UnexpectedChildError) Error() string {
	return fmt.Sprintf("parent: %s, child: %s", reflect.TypeOf(uc.parent).Name(), reflect.TypeOf(uc.child).Name())
}
