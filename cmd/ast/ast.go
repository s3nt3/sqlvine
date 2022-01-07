package main

import (
	"flag"
	"strings"

	"github.com/pingcap/tidb/parser/ast"

	"github.com/s3nt3/sqlvine/internal/logger"
	"github.com/s3nt3/sqlvine/pkg/parser"
)

var (
	opt_sql = flag.String("sql", "SELECT c1, c2 FROM t1;", "SQL to tst")
)

type visitor struct {
	level int
}

func (v *visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	defer func() { v.level++ }()

	if _, ok := in.(ast.StmtNode); ok {
		logger.L.Debugf("%s - %T %+v", strings.Repeat(" ", v.level), in, in)
	} else {
		switch in.(type) {
		case *ast.SubqueryExpr:
			logger.L.Debugf("%s - %T %+v", strings.Repeat(" ", v.level), in, in)
		case *ast.TableName:
			logger.L.Debugf("%s - %T %+v", strings.Repeat(" ", v.level), in, in)
		case *ast.TableSource:
			logger.L.Debugf("%s - %T %+v", strings.Repeat(" ", v.level), in, in)
		default:
			logger.L.Debugf("%s - %T", strings.Repeat(" ", v.level), in)
		}
	}

	return in, false
}

func (v *visitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	defer func() { v.level-- }()
	return in, true
}

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

	for _, stmt := range stmts {
		stmt.Accept(&visitor{})
	}
}
