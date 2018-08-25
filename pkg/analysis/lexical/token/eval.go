package token

import (
	"fmt"

	"github.com/bryanl/jsonnet-language-server/pkg/analysis/lexical/astext"
	"github.com/google/go-jsonnet/ast"
)

type evaluator struct {
	until ast.Node
	scope evalScope
}

type evalScope map[ast.Identifier]ast.Node

func (es evalScope) Clone() evalScope {
	clone := evalScope{}
	for k, v := range es {
		clone[k] = v
	}

	return clone
}

func (e *evaluator) eval(n ast.Node, parentScope evalScope) {
	switch n := n.(type) {
	case *ast.Local:
		s := parentScope.Clone()

		for _, bind := range n.Binds {
			s[bind.Variable] = bind.Body
		}

		e.eval(n.Body, s)
	case *astext.Partial:
		// nothing to do
	case *ast.Var:
		// nothing to do
	default:
		panic(fmt.Sprintf("unexpected node %T", n))
	}

	if n == e.until {
		e.scope = parentScope
	}
}

func eval(node, until ast.Node) evalScope {
	e := evaluator{until: until}

	s := evalScope{}
	e.eval(node, s)

	return e.scope
}