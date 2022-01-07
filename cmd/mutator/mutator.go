package main

import (
	"flag"

	"github.com/pingcap/tidb/parser/ast"
	"github.com/s3nt3/sqlvine/internal/logger"
	"github.com/s3nt3/sqlvine/pkg/mutator"
	"github.com/s3nt3/sqlvine/pkg/parser"
)

var (
	opt_sql   = flag.String("sql", "SELECT c1, c2 FROM t1;", "SQL to tst")
	opt_times = flag.Int("times", 3, "Mutate time")
)

func parse(sql string) []ast.StmtNode {
	parser := parser.NewTiDBParser()
	stmts, warns, err := parser.Parse([]byte(sql))
	if err != nil {
		logger.L.Panic(err.Error())
	}

	if len(warns) > 0 {
		for _, warn := range warns {
			logger.L.Debug(warn.Error())
		}
	}

	return stmts
}

func mutate(stmt ast.StmtNode) string {
	m := mutator.NewMutator(`[{
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

	stmt.Accept(m.Candidate)

	m.Mutate()

	sql, err := m.Candidate.GetTree().GetRoot().GetNode().Restore()
	if err != nil {
		logger.L.Panic(err.Error())
	}

	return sql
}

func main() {
	flag.Parse()

	sql := *opt_sql
	logger.L.Debugf("origin: %s", sql)

	stmts := parse(sql)
	for _, stmt := range stmts {
		for i := 0; i < *opt_times; i++ {
			sql = mutate(stmt)
			logger.L.Debugf("mutate[%.2d]: %s", i, sql)
		}
	}
}
