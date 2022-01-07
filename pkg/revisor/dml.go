package revisor

import (
	"fmt"
	"log"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	driver "github.com/pingcap/tidb/types/parser_driver"
	"github.com/s3nt3/sqlvine/internal/ir"
	"github.com/s3nt3/sqlvine/pkg/generator"
	"github.com/s3nt3/sqlvine/pkg/schema"
)

func (v *Revisor) reviseSelectStmt(node *ir.RevNode) {
	stmt := node.Node.(*ast.SelectStmt)
	if stmt.From == nil {
		v.reviseFrom(node)
	}

	if stmt.Where == nil {
		v.reviseWhere(node)
	}

	if stmt.Limit != nil {
		v.reviseLimit(node)
	}

	v.reviseFields(node)
}

func (v *Revisor) reviseFrom(node *ir.RevNode) {
	stmt := node.Node.(*ast.SelectStmt)
	g := generator.NewASTGenerator()
	stmt.From = g.TableRefsClause(v.schema.RandomNum(3))
}

func (v *Revisor) reviseWhere(node *ir.RevNode) {
	stmt := node.Node.(*ast.SelectStmt)
	g := generator.NewASTGenerator()
	stmt.Where = g.BinaryOperationExpr(v.schema.RandomNum(3))
}

func (v *Revisor) reviseLimit(node *ir.RevNode) {
	if constraint := node.GetParent(); constraint != nil {
		if expr, ok := constraint.Node.(*ast.SubqueryExpr); ok {
			if !expr.MultiRows {
				stmt := node.Node.(*ast.SelectStmt)
				stmt.Limit.Count.(*driver.ValueExpr).SetInt64(int64(1))
			}
		}
	}
}

func (v *Revisor) reviseFields(node *ir.RevNode) {
	stmt := node.Node.(*ast.SelectStmt)
	if len(stmt.Fields.Fields) == 0 {
		g := generator.NewASTGenerator()
		stmt.Fields.Fields = append(stmt.Fields.Fields, g.SelectField())
	}
}

func (v *Revisor) walkSelectStmt(node *ir.RevNode) {
	stmt := node.Node.(*ast.SelectStmt)
	if stmt.From != nil {
		v.walkFrom(node)
	}

	if stmt.Where != nil {
		v.walkWhere(node)
	}

	if stmt.GroupBy != nil {
		v.walkGroupBy(node)
	}

	if stmt.Having != nil {
		v.walkHaving(node)
	}

	if stmt.OrderBy != nil {
		v.walkOrderBy(node)
	}

	v.walkFields(node)
}

func (v *Revisor) walkFrom(node *ir.RevNode) *schema.Table {
	stmt := node.GetStmt()
	schema := stmt.GetSchema()

	if stmt.Walked {
		if table, ok := schema.TableMap["subquery"]; ok {
			return table
		}
	}

	from := node.GetChildByNodePtr(node.Node.(*ast.SelectStmt).From)
	table := v.walkTableRefsClause(from)
	schema.TableMap["subquery"] = table
	stmt.Walked = true

	return table
}

func (v *Revisor) walkTableRefsClause(node *ir.RevNode) *schema.Table {
	join := node.GetChildByNodePtr(node.Node.(*ast.TableRefsClause).TableRefs)
	return v.walkJoin(join)
}

func (v *Revisor) walkJoin(node *ir.RevNode) *schema.Table {
	join := node.Node.(*ast.Join)
	if join.Left != nil {
		if join.Right != nil {
			ltable := v.walkResultSetNode(node.GetChildByNodePtr(join.Left))
			rtable := v.walkResultSetNode(node.GetChildByNodePtr(join.Right))

			if join.On != nil {
				v.walkOnCondition(node.GetChildByNodePtr(join.On), []*schema.Table{ltable, rtable})
			}

			stmt := node.GetStmt()
			s := stmt.GetSchema()

			table := schema.MergeTable(ltable, rtable)
			table.SetID(int64(stmt.GetDepth()*10000 + len(v.schema.TableVec)*1000 + len(s.TableVec)))
			table.ReName(fmt.Sprintf("t%d", table.ID))

			return table
		} else {
			return v.walkResultSetNode(node.GetChildByNodePtr(join.Left))
		}
	}

	return nil
}

func (v *Revisor) walkOnCondition(node *ir.RevNode, tables []*schema.Table) {
	v.walkExprNode(node.GetChildByNodePtr(node.Node.(*ast.OnCondition).Expr), tables, nil)
}

func (v *Revisor) walkResultSetNode(node *ir.RevNode) *schema.Table {
	switch node.Node.(type) {
	case *ast.Join:
		return v.walkJoin(node)
	case *ast.SelectStmt:
		return v.walkFrom(node)
	case *ast.TableName:
		return v.walkTableName(node)
	case *ast.TableSource:
		return v.walkTableSource(node)
		// TODO: support *ast.SetOprStmt
		// case *ast.SetOprStmt:
	}
	return nil
}

func (v *Revisor) walkTableName(node *ir.RevNode) *schema.Table {
	tableName := node.Node.(*ast.TableName)
	if table, ok := v.schema.TableMap[tableName.Name.String()]; ok {
		return table
	}

	table := v.schema.GetRandomTable()
	if table != nil {
		tableName.Name = model.NewCIStr(table.Name)
	} else {
		log.Fatalf("Globle schema is empty: %+v", v.schema)
	}

	return table
}

func (v *Revisor) walkTableSource(node *ir.RevNode) *schema.Table {
	tableSource := node.Node.(*ast.TableSource)
	table := v.walkResultSetNode(node.GetChildByNodePtr(tableSource.Source))

	// if table not in global schema, add it into stmt schema as an inner table
	if _, ok := v.schema.TableMap[table.Name]; !ok {
		tableSource.AsName = model.NewCIStr(table.Name)
	}

	stmt := node.GetStmt()
	if !stmt.Walked {
		schema := stmt.GetSchema()
		schema.AddTable(table)
	}

	return table
}

func (v *Revisor) walkFields(node *ir.RevNode) {
	if constraint := node.GetParent(); constraint != nil {
		if _, ok := constraint.Node.(ast.ExprNode); ok {
			v.walkFieldList(node.GetChildByNodePtr(node.Node.(*ast.SelectStmt).Fields), false)
		}
	}

	v.walkFieldList(node.GetChildByNodePtr(node.Node.(*ast.SelectStmt).Fields), true)
}

func (v *Revisor) walkFieldList(node *ir.RevNode, multiCols bool) {
	for _, child := range node.GetChildren() {
		v.walkSelectField(child)
	}

	if !multiCols {
		fieldList := node.Node.(*ast.FieldList)
		if len(fieldList.Fields) != 1 {
			fieldList.Fields = []*ast.SelectField{
				fieldList.Fields[v.schema.RandomNum(len(fieldList.Fields))],
			}
		}
	}
}

func (v *Revisor) walkSelectField(node *ir.RevNode) {
	v.walkExprNode(node.GetChildByNodePtr(node.Node.(*ast.SelectField).Expr), []*schema.Table{}, nil)
}

func (v *Revisor) walkWhere(node *ir.RevNode) {
	v.walkExprNode(node.GetChildByNodePtr(node.Node.(*ast.SelectStmt).Where), []*schema.Table{}, nil)
}

func (v *Revisor) walkHaving(node *ir.RevNode) {
	having := node.GetChildByNodePtr(node.Node.(*ast.SelectStmt).Having)
	v.walkExprNode(having.GetChildByNodePtr(having.Node.(*ast.HavingClause).Expr), []*schema.Table{}, nil)
}

func (v *Revisor) walkOrderBy(node *ir.RevNode) {
	v.walkOrderByClause(node.GetChildByNodePtr(node.Node.(*ast.SelectStmt).OrderBy))
}

func (v *Revisor) walkOrderByClause(node *ir.RevNode) {
	for _, item := range node.GetChildren() {
		v.walkByItem(item)
	}
}

func (v *Revisor) walkGroupBy(node *ir.RevNode) {
	v.walkGroupByClause(node.GetChildByNodePtr(node.Node.(*ast.SelectStmt).GroupBy))
}

func (v *Revisor) walkGroupByClause(node *ir.RevNode) {
	for _, item := range node.GetChildren() {
		v.walkByItem(item)
	}
}

func (v *Revisor) walkByItem(node *ir.RevNode) {
	v.walkExprNode(node.GetChildByNodePtr(node.Node.(*ast.ByItem).Expr), []*schema.Table{}, nil)
}
