package revisor

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/s3nt3/sqlvine/internal/ir"
	"github.com/s3nt3/sqlvine/internal/session"
)

type Revisor struct {
	schema         *session.Schema
	tree           *ir.RevTree
	CandidateNodes []*ir.RevNode
	CandidateStmts []*ir.RevStmt
	err            error
}

func NewRevisor(schema *session.Schema) *Revisor {
	return &Revisor{
		schema:         schema,
		CandidateNodes: []*ir.RevNode{},
		CandidateStmts: []*ir.RevStmt{},
	}
}

func (v *Revisor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	if v.tree != nil {
		node := ir.NewRevNode(in)
		stmt := v.tree.GetCurrentStmt()
		node.SetParent(stmt.GetCurrentNode())

		if _, ok := in.(ast.StmtNode); ok {
			stmt := ir.NewRevStmt(node)
			stmt.SetParent(v.tree.GetCurrentStmt())
			v.tree.Push(stmt)
		} else {
			stmt.Push(node)
		}
	} else {
		v.tree = ir.NewRevTree(ir.NewRevStmt(ir.NewRevNode(in)))
	}

	v.revise(v.tree.GetCurrentStmt().GetCurrentNode())

	return in, v.err != nil
}

func (v *Revisor) Leave(in ast.Node) (out ast.Node, ok bool) {
	defer func() {
		v.tree.GetCurrentStmt().Pop()
		if _, ok := in.(ast.StmtNode); ok {
			v.tree.Pop()
		}
	}()

	v.walk(v.tree.GetCurrentStmt().GetCurrentNode())

	return in, v.err == nil
}

func (v *Revisor) revise(node *ir.RevNode) {
	if _, ok := node.Node.(ast.StmtNode); ok {
		v.reviseStmt(node)
	}
}

func (v *Revisor) reviseStmt(node *ir.RevNode) {
	switch node.Node.(type) {
	case *ast.SelectStmt:
		v.reviseSelectStmt(node)
	}
}

func (v *Revisor) walk(node *ir.RevNode) {
	if _, ok := node.Node.(ast.StmtNode); ok {
		v.walkStmt(node)
	}
}

func (v *Revisor) walkStmt(node *ir.RevNode) {
	switch node.Node.(type) {
	case *ast.SelectStmt:
		v.walkSelectStmt(node)
	}
}
