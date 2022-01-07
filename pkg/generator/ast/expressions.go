package ast

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/opcode"
	driver "github.com/pingcap/tidb/types/parser_driver"
)

func (g *ASTGenerator) BinaryOperationExpr(depth int) ast.ExprNode {
	expr := &ast.BinaryOperationExpr{}
	if depth == 0 {
		switch g.RandomNum(3) {
		case 0:
			expr.Op = opcode.GT
		case 1:
			expr.Op = opcode.LT
		case 2:
			expr.Op = opcode.NE
		case 3:
			expr.Op = opcode.EQ
		}

		if g.RandomBool() {
			expr.L = g.ExprNode(true)
			expr.R = g.ExprNode(false)
		} else {
			expr.L = g.ColumnNameExpr()
			expr.R = g.ValueExpr()
		}
	} else {
		switch g.RandomNum(3) {
		case 0:
			expr.Op = opcode.LogicAnd
		case 1:
			expr.Op = opcode.LogicOr
		case 2:
			expr.Op = opcode.LogicXor
		}

		expr.L = g.BinaryOperationExpr(0)
		expr.R = g.BinaryOperationExpr(depth - 1)
	}

	return expr
}

func (g *ASTGenerator) ColumnNameExpr() ast.ExprNode {
	return &ast.ColumnNameExpr{}
}

func (g *ASTGenerator) SubqueryExpr() ast.ExprNode {
	multiRows := g.RandomBool()

	return &ast.SubqueryExpr{
		Query:     g.SelectStmt(1, multiRows),
		MultiRows: multiRows,
	}
}

func (g *ASTGenerator) ValueExpr() ast.ExprNode {
	return &driver.ValueExpr{}
}

func (g *ASTGenerator) ExprNode(column bool) ast.ExprNode {
	if g.RandomNum(10) > 0 {
		if column {
			return g.ColumnNameExpr()
		} else {
			switch g.RandomNum(3) {
			case 0:
				return g.ValueExpr()
			default:
				return g.ColumnNameExpr()
			}
		}
	} else {
		return g.SubqueryExpr()
	}
}
