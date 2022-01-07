package ast

import "github.com/pingcap/tidb/parser/ast"

func (g *ASTGenerator) AggregateFuncExpr() ast.ExprNode {
	candidate := [...]string{
		"count",
		"sum",
		"avg",
		"max",
		"min",
		"group_concat",
		"var_pop",
		"var_samp",
		"stddev_pop",
		"stddev_samp",
		"json_objectagg",
	}

	switch funcName := candidate[g.RandomNum(len(candidate))]; funcName {
	case "json_objectagg":
		return &ast.AggregateFuncExpr{
			F: funcName,
			Args: []ast.ExprNode{
				g.ExprNode(true),
				g.ExprNode(true),
			},
		}
	default:
		return &ast.AggregateFuncExpr{
			F: funcName,
			Args: []ast.ExprNode{
				g.ExprNode(true),
			},
		}
	}
}
