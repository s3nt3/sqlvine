package session

import (
	"sync"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"

	_ "github.com/pingcap/tidb/types/parser_driver"
)

type TiDBParser struct {
	parser *parser.Parser
	mutex  *sync.Mutex
}

func NewTiDBParser() *TiDBParser {
	return &TiDBParser{
		parser: parser.New(),
		mutex:  &sync.Mutex{},
	}
}

func (p *TiDBParser) Parse(bytes []byte) ([]ast.StmtNode, []error, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// TODO: need to support configurable charsets and collation.
	// We send empty string as the 2nd and 3rd arguments of parser.Parse() to
	// use the default charset and collation right now.
	return p.parser.Parse(string(bytes), "", "")
}
