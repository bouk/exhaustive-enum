package main

import (
	"go/ast"
)

func getInnerFromParen(e ast.Expr) ast.Expr {
	for e2, ok := e.(*ast.ParenExpr); ok; {
		e = e2.X
	}
	return e
}
