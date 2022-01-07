package mutator

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/s3nt3/sqlvine/internal/ir"
	"github.com/s3nt3/sqlvine/internal/util"
	"github.com/s3nt3/sqlvine/pkg/parser"
	"github.com/s3nt3/sqlvine/pkg/revisor"
	"github.com/s3nt3/sqlvine/pkg/schema"
)

type Mutator struct {
	Mutated int

	Candidate *MutationCandidate
	Parser    *parser.TiDBParser
	Revisor   *revisor.Revisor
	Schema    *schema.Schema

	*util.Random
}

func NewMutator(s string) *Mutator {
	schema := schema.NewSchema(s)
	return &Mutator{
		Candidate: NewMutationCandidate(),
		Random:    util.NewRandom(),
		Parser:    parser.NewTiDBParser(),
		Revisor:   revisor.NewRevisor(schema),
		Schema:    schema,
	}
}

func (m *Mutator) updateMutateStatus(node *ir.MutNode) {
	for tmp := node; tmp != nil && tmp.Mutated != true; tmp = tmp.GetParent() {
		m.Candidate.GetTree().Mutated++
		tmp.Mutated = true
	}
	m.Mutated++
}

func (m *Mutator) randomCandidateNode() *ir.MutNode {
	for len(m.Candidate.CandidateNodes) > 0 && m.Mutated < len(m.Candidate.CandidateNodes) {
		idx := m.RandomNum(len(m.Candidate.CandidateNodes))
		candidate := m.Candidate.CandidateNodes[idx]
		if !candidate.Mutated {
			return candidate
		} else {
			m.Candidate.CandidateNodes = append(m.Candidate.CandidateNodes[:idx], m.Candidate.CandidateNodes[idx+1:]...)
		}
	}

	return nil
}

func (m *Mutator) Mutate() {
	for m.Candidate.GetTree().Mutated < m.Candidate.GetTree().Total/2 {
		if node := m.randomCandidateNode(); node != nil {
			m.MutateNode(node)
		} else {
			break
		}
	}

	for _, stmt := range m.Candidate.CandidateStmts {
		if !stmt.GetNode().Mutated {
			m.MutateNode(stmt.GetNode())
		}
	}

	m.Candidate.GetTree().GetRoot().GetNode().Node.Accept(m.Revisor)
}

func (m *Mutator) MutateNode(node *ir.MutNode) {
	defer m.updateMutateStatus(node)

	switch node.GetStmt().GetNode().Node.(type) {
	case *ast.SelectStmt:
		m.mutateSelectStmt(node)
	}
}
