package mutator

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/s3nt3/sqlvine/internal/ir"
	"github.com/s3nt3/sqlvine/pkg/generator"
)

func (m *Mutator) mutateSelectStmt(node *ir.MutNode) {
	switch node.Node.(type) {
	case *ast.SelectStmt:
		m.MutateSelectStmtNode(node)
	case *ast.FieldList:
		m.MutateFieldList(node)
	case *ast.SelectField:
		m.MutateSelectField(node)
	}
}

func (m *Mutator) MutateSelectStmtNode(node *ir.MutNode) {
	stmt := node.Node.(*ast.SelectStmt)

	switch {
	case stmt.From == nil:
		stmt.From = generator.NewASTGenerator().TableRefsClause(m.RandomNum(1))
	case stmt.Where == nil:
		stmt.Where = generator.NewASTGenerator().BinaryOperationExpr(m.RandomNum(1))
	}

	switch m.RandomNum(6) {
	case 0:
		stmt.From = generator.NewASTGenerator().TableRefsClause(m.RandomNum(3))
	case 1:
		stmt.Where = generator.NewASTGenerator().BinaryOperationExpr(m.RandomNum(3))
	case 2:
		stmt.GroupBy = generator.NewASTGenerator().GroupByClause()
	case 3:
		if stmt.GroupBy == nil {
			stmt.GroupBy = generator.NewASTGenerator().GroupByClause()
		}
		stmt.Having = generator.NewASTGenerator().HavingClause(m.RandomNum(3))
	case 4:
		stmt.OrderBy = generator.NewASTGenerator().OrderByClause()
	case 5:
		stmt.Limit = generator.NewASTGenerator().Limit(0, 1)
	default:
		switch {
		case stmt.Limit != nil:
			stmt.Limit = nil
		case stmt.OrderBy != nil:
			stmt.OrderBy = nil
		case stmt.Having != nil:
			stmt.Having = nil
		}
	}
}

func (m *Mutator) MutateFieldList(node *ir.MutNode) {
	list := node.Node.(*ast.FieldList)

	switch m.RandomNum(2) {
	case 0:
		len := m.RandomNum(16/len(list.Fields) + 1)
		for i := 0; i < len; i++ {
			list.Fields = append(list.Fields, generator.NewASTGenerator().SelectField())
		}
	case 1:
		idx := m.RandomNum(len(list.Fields))
		if len(list.Fields) > 0 {
			list.Fields = append(list.Fields[:idx], list.Fields[idx+1:]...)
		}
	}
}

func (m *Mutator) MutateSelectField(node *ir.MutNode) {
	field := node.Node.(*ast.SelectField)
	stmt := node.GetStmt().GetNode().Node.(*ast.SelectStmt)

	if stmt.GroupBy != nil {
		if _, ok := field.Expr.(*ast.ColumnNameExpr); ok {
			switch m.RandomNum(3) {
			case 0:
				field.Expr = generator.NewASTGenerator().ExprNode(true)
			default:
				field.Expr = generator.NewASTGenerator().AggregateFuncExpr()
			}
		}
	} else {
		field.Expr = generator.NewASTGenerator().ExprNode(true)
	}
}
