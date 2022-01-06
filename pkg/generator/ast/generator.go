package ast

import (
	"github.com/pingcap/tidb/parser/ast"

	"github.com/s3nt3/sqlvine/internal/util"
)

type ASTGenerator struct {
	*util.Random

	Node ast.Node
}
