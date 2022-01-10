package plan_test

import (
	"testing"

	"plan"

	"github.com/s3nt3/sqlvine/pkg/mutator"
	"github.com/s3nt3/sqlvine/pkg/parser"
	"github.com/s3nt3/sqlvine/pkg/schema"
)

type SQLMutator struct {
	sql     []byte
	payload []byte
}

func NewSQLMutator(sql string) *SQLMutator {
	return &SQLMutator{
		sql:     []byte(sql),
		payload: []byte(sql),
	}
}

func (m *SQLMutator) Payload() []byte {
	return m.payload
}

func (m *SQLMutator) Unmarshal(sql []byte) error {
	m.sql = sql
	return nil
}

func (m *SQLMutator) Marshal() ([]byte, error) {
	return m.sql, nil
}

func (m *SQLMutator) Mutate() error {
	stmts, _, err := parser.NewTiDBParser().Parse(m.sql)
	if err != nil {
		return nil
	}

	s := schema.NewSchema(`[{
		"id": 1,
		"name": "t1",
		"charset": "utf8mb4",
		"collate": "utf8mb4_bin",
	
		"columns": [{
			"id": 1,
			"name": "c1",
			"table": "t1",
			"type": "int",
			"primary_key": true
		},{
			"id": 2,
			"name": "c2",
			"table": "t1",
			"type": "varchar",
			"size": 100
		},{
			"id": 3,
			"name": "c3",
			"table": "t1",
			"type": "varchar",
			"size": 100
		}],
		"indices": []
	},{
		"id": 2,
		"name": "t2",
		"charset": "utf8mb4",
		"collate": "utf8mb4_bin",
	
		"columns": [{
			"id": 1,
			"name": "c1",
			"table": "t1",
			"type": "int",
			"primary_key": true
		},{
			"id": 2,
			"name": "c2",
			"table": "t1",
			"type": "varchar",
			"size": 100
		},{
			"id": 3,
			"name": "c3",
			"table": "t1",
			"type": "varchar",
			"size": 100
		}],
		"indices": []
	},{
		"id": 3,
		"name": "t3",
		"charset": "utf8mb4",
		"collate": "utf8mb4_bin",
	
		"columns": [{
			"id": 1,
			"name": "c1",
			"table": "t1",
			"type": "int",
			"primary_key": true
		},{
			"id": 2,
			"name": "c2",
			"table": "t1",
			"type": "varchar",
			"size": 100
		},{
			"id": 3,
			"name": "c3",
			"table": "t1",
			"type": "varchar",
			"size": 100
		}],
		"indices": []
	
	}]`)

	mut := mutator.NewMutator(s)

	stmts[0].Accept(mut.Candidate)

	mut.Mutate()

	sql, err := mut.Candidate.GetTree().GetRoot().GetNode().Restore()
	if err != nil {
		return err
	}
	m.payload = []byte(sql)

	return nil
}

func FuzzBuildLogicalPlan(f *testing.F) {
	m := NewSQLMutator("SELECT c1 FROM t1;")

	f.Add(m)

	f.Fuzz(func(t *testing.T, m *SQLMutator) {
		// plan.BuildLogicalPlan(string(m.Payload()))
		plan.PrintSQL(string(m.Payload()))
	})
}
