package plan

import (
	"log"
	"strings"

	"github.com/pingcap/tidb/planner/core"
	"github.com/s3nt3/sqlvine/pkg/parser"
	"github.com/s3nt3/sqlvine/pkg/planner"
	"github.com/s3nt3/sqlvine/pkg/schema"
)

func BuildLogicalPlan(sql string) {
	log.Printf("[debug] parse sql: %s", sql)

	parser := parser.NewTiDBParser()
	_, warns, err := parser.Parse([]byte(sql))
	if err != nil {
		log.Printf("[debug] parse failed: %s", err.Error())
		return
	}

	if len(warns) > 0 {
		for _, warn := range warns {
			log.Printf("[debug] parse warning: %s", warn.Error())
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

	_ = planner.NewTiDBPlanBuilder(schema.GetSchemaInfo())
	// for _, stmt := range stmts {
	// 	plan, _, err := builder.Build(stmt)
	// 	if err != nil {
	// 		log.Printf("[debug] plan build failed: %s", err.Error())
	// 	} else {
	// 		PrintLogicalPlan(plan.(core.LogicalPlan), 0)
	// 	}
	// }
}

func PrintLogicalPlan(plan core.LogicalPlan, level int) {
	if plan != nil {
		switch n := plan.(type) {
		case *core.DataSource:
			log.Printf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalAggregation:
			log.Printf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalJoin:
			log.Printf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalLimit:
			log.Printf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalProjection:
			log.Printf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalSelection:
			log.Printf("%s|-%T", strings.Repeat("\t", level), plan)
		case *core.LogicalSort:
			log.Printf("%s|-%T", strings.Repeat("\t", level), plan)
		default:
			log.Printf("[error] `%T` not supported", n)
		}

		for _, p := range plan.Children() {
			PrintLogicalPlan(p, level+1)
		}
	}
}

func PrintSQL(sql string) {
	log.Println(sql)
}
