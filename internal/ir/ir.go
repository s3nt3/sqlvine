package ir

import (
	"bytes"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/format"
	"github.com/s3nt3/sqlvine/internal/session"
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
	ctx := format.NewRestoreCtx(format.RestoreKeyWordUppercase|format.RestoreNameLowercase|format.RestoreNameBackQuotes|format.RestoreStringWithoutDefaultCharset, buf)
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

type RevNode struct {
	Node     ast.Node
	parent   *RevNode
	children []*RevNode
	stmt     *RevStmt
	depth    int
}

func NewRevNode(node ast.Node) *RevNode {
	return &RevNode{
		Node: node,
	}
}

func (n *RevNode) addChild(child *RevNode) {
	if child != nil {
		n.children = append(n.children, child)
	}
}

func (n *RevNode) GetChildren() []*RevNode {
	return n.children
}

func (n *RevNode) GetChildByNodePtr(node ast.Node) *RevNode {
	for _, child := range n.children {
		switch c := child.Node.(type) {
		default:
			switch n := node.(type) {
			default:
				if c == n {
					return child
				}
			}
		}
	}

	return nil
}

func (n *RevNode) GetParent() *RevNode {
	return n.parent
}

func (n *RevNode) SetParent(parent *RevNode) {
	if parent != nil {
		n.parent = parent
		n.parent.addChild(n)
	}
}

func (n *RevNode) GetDepth() int {
	return n.depth
}

func (n *RevNode) setDepth(depth int) {
	n.depth = depth
}

func (n *RevNode) GetStmt() *RevStmt {
	return n.stmt
}

func (n *RevNode) setStmt(stmt *RevStmt) {
	if stmt != nil {
		n.stmt = stmt
		n.stmt.addMember(n)
	}
}

type RevStmt struct {
	node     *RevNode
	path     []*RevNode
	parent   *RevStmt
	children []*RevStmt
	members  []*RevNode
	schema   *session.Schema
	Walked   bool
	depth    int
}

func NewRevStmt(node *RevNode) *RevStmt {
	stmt := &RevStmt{
		node:   node,
		path:   []*RevNode{},
		schema: session.NewSchema("[]"),
	}
	stmt.Push(node)

	return stmt
}

func (n *RevStmt) addChild(child *RevStmt) {
	if child != nil {
		n.children = append(n.children, child)
	}
}

func (n *RevStmt) addMember(member *RevNode) {
	if member != nil {
		n.members = append(n.members, member)
	}
}

func (s *RevStmt) GetSchema() *session.Schema {
	return s.schema
}

func (s *RevStmt) GetParent() *RevStmt {
	return s.parent
}

func (s *RevStmt) SetParent(parent *RevStmt) {
	if parent != nil {
		s.parent = parent
		s.parent.addChild(s)
	}
}

func (s *RevStmt) GetDepth() int {
	return s.depth
}

func (s *RevStmt) setDepth(depth int) {
	s.depth = depth
}

func (s *RevStmt) GetPath() []*RevNode {
	return s.path
}

func (s *RevStmt) GetNode() *RevNode {
	return s.node
}

func (s *RevStmt) GetCurrentNode() *RevNode {
	if len(s.path) > 0 {
		return s.path[len(s.path)-1]
	}
	return nil
}

func (s *RevStmt) Push(node *RevNode) {
	if node != nil {
		node.setDepth(len(s.path))
		s.path = append(s.path, node)
		node.setStmt(s)
	}
}

func (s *RevStmt) Pop() *RevNode {
	node := s.GetCurrentNode()
	if node != nil {
		s.path = s.path[:len(s.path)-1]
	}
	return node
}

type RevTree struct {
	root *RevStmt
	path []*RevStmt
}

func NewRevTree(stmt *RevStmt) *RevTree {
	tree := &RevTree{
		root: stmt,
		path: []*RevStmt{},
	}
	tree.Push(stmt)

	return tree
}

func (t *RevTree) GetPath() []*RevStmt {
	return t.path
}

func (t *RevTree) GetRoot() *RevStmt {
	return t.root
}

func (t *RevTree) GetCurrentStmt() *RevStmt {
	if len(t.path) > 0 {
		return t.path[len(t.path)-1]
	}
	return nil
}

func (t *RevTree) Push(stmt *RevStmt) {
	if stmt != nil {
		stmt.setDepth(len(t.path))
		t.path = append(t.path, stmt)
	}
}

func (t *RevTree) Pop() *RevStmt {
	stmt := t.GetCurrentStmt()
	if stmt != nil {
		t.path = t.path[:len(t.path)-1]
	}
	return stmt
}
