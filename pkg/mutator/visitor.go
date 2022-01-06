package mutator

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/s3nt3/sqlvine/internal/ir"
)

type MutationCandidate struct {
	CandidateNodes []*ir.MutNode
	CandidateStmts []*ir.MutStmt

	err  error
	tree *ir.MutTree
}

func NewMutationCandidate() *MutationCandidate {
	return &MutationCandidate{
		CandidateNodes: []*ir.MutNode{},
		CandidateStmts: []*ir.MutStmt{},
	}
}

func (v *MutationCandidate) GetTree() *ir.MutTree {
	return v.tree
}

func (v *MutationCandidate) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	if v.tree != nil {
		node := ir.NewMutNode(in)
		stmt := v.tree.GetCurrentStmt()
		node.SetParent(stmt.GetCurrentNode())

		if _, ok := in.(ast.StmtNode); ok {
			stmt := ir.NewMutStmt(node)
			stmt.SetParent(v.tree.GetCurrentStmt())
			v.tree.Push(stmt)
		} else {
			stmt.Push(node)
		}
	} else {
		v.tree = ir.NewMutTree(ir.NewMutStmt(ir.NewMutNode(in)))
	}

	return in, v.err != nil
}

func (v *MutationCandidate) Leave(in ast.Node) (out ast.Node, ok bool) {
	defer func() {
		v.tree.GetCurrentStmt().Pop()
		if _, ok := in.(ast.StmtNode); ok {
			v.tree.Pop()
		}
		v.tree.Total++
	}()

	v.addCandidate(v.tree.GetCurrentStmt().GetCurrentNode())

	return in, v.err == nil
}

func (v *MutationCandidate) addCandidate(node *ir.MutNode) {
	switch node.Node.(type) {
	case *ast.SelectStmt:
		v.addCandidateNode(node)
	}

	if _, ok := node.Node.(ast.StmtNode); ok {
		v.addCandidateStmt(node)
	}
}

func (v *MutationCandidate) addCandidateNode(node *ir.MutNode) {
	v.CandidateNodes = append(v.CandidateNodes, node)
}

func (v *MutationCandidate) addCandidateStmt(node *ir.MutNode) {
	v.CandidateStmts = append(v.CandidateStmts, node.GetStmt())
}
