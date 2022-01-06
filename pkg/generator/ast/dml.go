package ast

import (
	"github.com/pingcap/tidb/parser/ast"

	driver "github.com/pingcap/tidb/types/parser_driver"
)

func (g *ASTGenerator) SelectStmt(depth int, multiRows bool) *ast.SelectStmt {
	selectStmt := ast.SelectStmt{
		SelectStmtOpts: &ast.SelectStmtOpts{
			SQLCache: true,
		},

		Fields: g.FieldList(),

		LockInfo: &ast.SelectLockInfo{},
	}

	if depth > 1 {
		depth = g.RandomNum(depth)
	} else {
		depth = 0
	}

	selectStmt.From = g.TableRefsClause(depth)
	selectStmt.Where = g.BinaryOperationExpr(depth)

	if !multiRows {
		selectStmt.Limit = g.Limit(0, 1)
	}

	return &selectStmt
}

func (g *ASTGenerator) FieldList() *ast.FieldList {
	return &ast.FieldList{
		Fields: []*ast.SelectField{
			g.SelectField(),
		},
	}
}

func (g *ASTGenerator) SelectField() *ast.SelectField {
	return &ast.SelectField{
		Expr: g.ColumnNameExpr(),
	}
}

func (g *ASTGenerator) TableRefsClause(depth int) *ast.TableRefsClause {
	tableRefsClause := &ast.TableRefsClause{
		TableRefs: g.Join(),
	}

	if depth > 0 {
		if g.RandomNum(3) > 0 {
			tableRefsClause.TableRefs.On = g.OnCondition(depth)
		}
		tableRefsClause.TableRefs.Right = g.TableSource(depth)
	}

	return tableRefsClause
}

func (g *ASTGenerator) Join() *ast.Join {
	return &ast.Join{
		Left: g.TableName(),
	}
}

func (g *ASTGenerator) OnCondition(depth int) *ast.OnCondition {
	return &ast.OnCondition{
		Expr: g.BinaryOperationExpr(depth - 1),
	}
}

func (g *ASTGenerator) TableSource(depth int) *ast.TableSource {
	if depth > 0 {
		return &ast.TableSource{
			Source: g.SelectStmt(depth-1, true),
		}
	} else {
		return &ast.TableSource{
			Source: g.TableName(),
		}
	}
}

func (g *ASTGenerator) TableName() *ast.TableName {
	return &ast.TableName{}
}

func (g *ASTGenerator) GroupByClause() *ast.GroupByClause {
	return &ast.GroupByClause{
		Items: []*ast.ByItem{
			{
				Expr: g.ExprNode(true),
			},
		},
	}
}

func (g *ASTGenerator) HavingClause(depth int) *ast.HavingClause {
	return &ast.HavingClause{
		Expr: g.BinaryOperationExpr(depth),
	}
}

func (g *ASTGenerator) OrderByClause() *ast.OrderByClause {
	return &ast.OrderByClause{
		Items: []*ast.ByItem{
			{
				Expr: g.ExprNode(true),
			},
		},
	}
}

func (g *ASTGenerator) Limit(offset int, count int) *ast.Limit {
	c := &driver.ValueExpr{}
	c.SetInt64(int64(count))

	if offset > 0 {
		o := &driver.ValueExpr{}
		o.SetInt64(int64(offset))

		return &ast.Limit{
			Count:  c,
			Offset: o,
		}
	}

	return &ast.Limit{
		Count: c,
	}
}
