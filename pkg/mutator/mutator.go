package mutator

import (
	"github.com/pingcap/tidb/parser/ast"
	"github.com/s3nt3/sqlvine/internal/ir"
	"github.com/s3nt3/sqlvine/internal/session"
	"github.com/s3nt3/sqlvine/internal/util"
	"github.com/s3nt3/sqlvine/pkg/revisor"
)

type Mutator struct {
	Candidate *MutationCandidate

	Mutated int

	revisor *revisor.Revisor

	*util.Random
}

func NewMutator(schema *session.Schema) *Mutator {
	return &Mutator{
		Candidate: NewMutationCandidate(),
		Random:    util.NewRandom(),
		revisor:   revisor.NewRevisor(schema),
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

	m.Candidate.GetTree().GetRoot().GetNode().Node.Accept(m.revisor)
}

func (m *Mutator) MutateNode(node *ir.MutNode) {
	defer m.updateMutateStatus(node)

	switch node.GetStmt().GetNode().Node.(type) {
	case *ast.SelectStmt:
		m.mutateSelectStmt(node)
	}
}
