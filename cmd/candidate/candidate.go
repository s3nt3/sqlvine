package main

import (
	"flag"
	"strings"

	"github.com/s3nt3/sqlvine/internal/logger"
	"github.com/s3nt3/sqlvine/internal/session"
	"github.com/s3nt3/sqlvine/pkg/mutator"
)

var (
	opt_sql = flag.String("sql", "SELECT c1, c2 FROM t1;", "SQL to tst")
)

func main() {
	flag.Parse()

	parser := session.NewTiDBParser()
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
		candidate := mutator.NewMutationCandidate()
		stmt.Accept(candidate)

		for idx, node := range candidate.CandidateNodes {
			sql, err := node.Restore()
			if err != nil {
				logger.L.Panic(err.Error())
			} else {
				logger.L.Debugf("node[%.2d] %+v", idx, node)
				logger.L.Debugf("\t%s - %s", strings.Repeat("\t", node.GetDepth()), sql)
			}
		}

		for idx, stmt := range candidate.CandidateStmts {
			sql, err := stmt.GetNode().Restore()
			if err != nil {
				logger.L.Panic(err.Error())
			} else {
				logger.L.Debugf("stmt[%.2d] %+v", idx, stmt)
				logger.L.Debugf("\t%s - %s", strings.Repeat("\t", stmt.GetDepth()), sql)
			}
		}
	}
}
