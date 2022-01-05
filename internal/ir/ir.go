package ir

import (
	"bytes"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
)

type MutNode struct {
	Node    ast.Node
	Mutated bool

	parent *MutNode
	stmt   *MutStmt
	depth  int
}

func NewMutNode(node ast.Node) *MutNode {
	return &MutNode{
		Node: node,
	}
}

func (n *MutNode) GetParent() *MutNode {
	return n.parent
}

func (n *MutNode) SetParent(parent *MutNode) {
	if parent != nil {
		n.parent = parent
	}
}

func (n *MutNode) GetStmt() *MutStmt {
	return n.stmt
}

func (n *MutNode) setStmt(stmt *MutStmt) {
	if stmt != nil {
		n.stmt = stmt
	}
}

func (n *MutNode) GetDepth() int {
	return n.depth
}

func (n *MutNode) setDepth(depth int) {
	n.depth = depth
}

func (n *MutNode) Restore() (string, error) {
	buf := new(bytes.Buffer)
	ctx := format.NewRestoreCtx(format.RestoreKeyWordUppercase|format.RestoreNameLowercase|format.RestoreNameBackQuotes, buf)
	err := n.Node.Restore(ctx)
	if nil != err {
		return "", err
	}
	return buf.String(), nil
}

type MutStmt struct {
	node *MutNode
	path []*MutNode

	parent *MutStmt
	depth  int
}

func NewMutStmt(node *MutNode) *MutStmt {
	stmt := &MutStmt{
		node: node,
		path: []*MutNode{},
	}
	stmt.Push(node)

	return stmt
}

func (s *MutStmt) GetParent() *MutStmt {
	return s.parent
}

func (s *MutStmt) SetParent(parent *MutStmt) {
	if parent != nil {
		s.parent = parent
	}
}

func (s *MutStmt) GetPath() []*MutNode {
	return s.path
}

func (s *MutStmt) GetNode() *MutNode {
	return s.node
}

func (s *MutStmt) GetDepth() int {
	return s.depth
}

func (s *MutStmt) setDepth(depth int) {
	s.depth = depth
}

func (s *MutStmt) GetCurrentNode() *MutNode {
	if len(s.path) > 0 {
		return s.path[len(s.path)-1]
	}
	return nil
}

func (s *MutStmt) Push(node *MutNode) {
	if node != nil {
		node.setDepth(len(s.path))
		s.path = append(s.path, node)
		node.setStmt(s)
	}
}

func (s *MutStmt) Pop() *MutNode {
	node := s.GetCurrentNode()
	if node != nil {
		s.path = s.path[:len(s.path)-1]
	}
	return node
}

type MutTree struct {
	root *MutStmt
	path []*MutStmt

	Mutated int
	Total   int
}

func NewMutTree(stmt *MutStmt) *MutTree {
	tree := &MutTree{
		root:  stmt,
		path:  []*MutStmt{},
		Total: 1,
	}
	tree.Push(stmt)

	return tree
}

func (t *MutTree) GetPath() []*MutStmt {
	return t.path
}

func (t *MutTree) GetRoot() *MutStmt {
	return t.root
}

func (t *MutTree) GetCurrentStmt() *MutStmt {
	if len(t.path) > 0 {
		return t.path[len(t.path)-1]
	}
	return nil
}

func (t *MutTree) Push(stmt *MutStmt) {
	if stmt != nil {
		stmt.setDepth(len(t.path))
		t.path = append(t.path, stmt)
	}
}

func (t *MutTree) Pop() *MutStmt {
	stmt := t.GetCurrentStmt()
	if stmt != nil {
		t.path = t.path[:len(t.path)-1]
	}
	return stmt
}
