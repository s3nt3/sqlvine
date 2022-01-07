package main

import (
	"flag"
	"strings"

	"github.com/pingcap/tidb/planner/core"
	"github.com/s3nt3/sqlvine/internal/session"
	"github.com/s3nt3/sqlvine/pkg/logger"
	"github.com/s3nt3/sqlvine/pkg/parser"
	"github.com/s3nt3/sqlvine/pkg/schema"
)

var (
	opt_sql = flag.String("sql", "SELECT c1, c2 FROM t1;", "SQL to build")
)

func main() {
	flag.Parse()

	parser := parser.NewTiDBParser()
	stmts, warns, err := parser.Parse([]byte(*opt_sql))
	if err != nil {
		logger.L.Panic(err.Error())
	}

	if len(warns) > 0 {
		for _, warn := range warns {
			logger.L.Debug(warn.Error())
		}
	}

	schema := schema.NewSchema(`[{
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
	builder := session.NewTiDBPlanBuilder(schema.GetSchemaInfo())
	for _, stmt := range stmts {
		plan, _, err := builder.Build(stmt)
		if err != nil {
			logger.L.Debug(err.Error())
		} else {
			PrintLogicalPlan(plan.(core.LogicalPlan), 0)
		}
	}
}

func PrintLogicalPlan(plan core.LogicalPlan, level int) {
	if plan != nil {
		switch n := plan.(type) {
		case *core.DataSource:
			logger.L.Debugf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalAggregation:
			logger.L.Debugf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalJoin:
			logger.L.Debugf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalLimit:
			logger.L.Debugf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalProjection:
			logger.L.Debugf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalSelection:
			logger.L.Debugf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalSort:
			logger.L.Debugf("%s|-%T", strings.Repeat("\t", level), plan)
		default:
			logger.L.Panicf("`%T` not supported", n)
		}

		for _, p := range plan.Children() {
			PrintLogicalPlan(p, level+1)
		}
	}
}
