package revisor

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/types"

	driver "github.com/pingcap/tidb/types/parser_driver"

	"github.com/s3nt3/sqlvine/internal/ir"
	"github.com/s3nt3/sqlvine/internal/logger"
	"github.com/s3nt3/sqlvine/internal/session"
	"github.com/s3nt3/sqlvine/pkg/generator"
)

func (v *Revisor) walkExprNode(node *ir.RevNode, tables []*session.Table, column *session.Column) *session.Column {
	switch node.Node.(type) {
	case *ast.BinaryOperationExpr:
		v.walkBinaryOperationExpr(node, tables)
	case *ast.ColumnNameExpr:
		return v.walkColumnNameExpr(node, tables)
	case *ast.SubqueryExpr:
		return v.walkSubqueryExpr(node).GetRandomColumn()
	case *driver.ValueExpr:
		v.walkValueExpr(node, column)
	default:
		logger.L.Panicf("Expr type `%T` not supported", node.Node)
	}

	return column
}

func (v *Revisor) walkBinaryOperationExpr(node *ir.RevNode, tables []*session.Table) {
	var table *session.Table

	if len(tables) > 0 {
		table = tables[v.schema.RandomNum(len(tables))]
	} else {
		table = node.GetStmt().GetSchema().GetRandomTable()
	}

	if table == nil {
		table = v.schema.GetRandomTable()
	}

	column := table.GetRandomColumn()
	expr := node.Node.(*ast.BinaryOperationExpr)
	v.walkExprNode(node.GetChildByNodePtr(expr.L), tables, v.walkExprNode(
		node.GetChildByNodePtr(expr.R), tables, column,
	))
}

func (v *Revisor) walkColumnNameExpr(node *ir.RevNode, tables []*session.Table) *session.Column {
	var table *session.Table

	if len(tables) > 0 {
		table = tables[v.schema.RandomNum(len(tables))]
	} else {
		table = node.GetStmt().GetSchema().GetRandomTable()
	}

	if table == nil {
		table = v.schema.GetRandomTable()
	}

	column := table.GetRandomColumn()

	expr := node.Node.(*ast.ColumnNameExpr)
	expr.Name = &ast.ColumnName{
		Table: model.NewCIStr(column.Table),
		Name:  model.NewCIStr(column.Name),
	}

	return column
}

func (v *Revisor) walkSubqueryExpr(node *ir.RevNode) *session.Table {
	expr := node.Node.(*ast.SubqueryExpr)
	switch expr.Query.(type) {
	case *ast.SelectStmt:
		return v.walkFrom(node.GetChildByNodePtr(expr.Query))
	default:
		logger.L.Panicf("Stmt `%T` not supported in SubqueryExpr", node.Node)
	}

	return nil
}

func (v *Revisor) walkValueExpr(node *ir.RevNode, column *session.Column) {
	expr := node.Node.(*driver.ValueExpr)
	if column != nil {
		g := generator.NewValueGenerator()
		switch column.Type {
		case "int":
			expr.SetInt64(int64(g.NewInt()))
		case "float":
			expr.SetFloat64(g.NewFloat())
		case "timestamp":
			expr.SetMysqlTime(types.NewTime(types.FromGoTime(g.NewTimestamp()), 0, 0))
		case "varchar":
			expr.SetString(g.NewString(), "utf8mb4_bin")
			expr.TexprNode.Type.Charset = "utf8mb4"
			expr.TexprNode.Type.Collate = "utf8mb4_bin"
		}
	}
}
