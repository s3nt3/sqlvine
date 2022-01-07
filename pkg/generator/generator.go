package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/s3nt3/sqlvine/internal/util"
	"github.com/s3nt3/sqlvine/pkg/generator/ast"
)

func NewASTGenerator() *ast.ASTGenerator {
	return &ast.ASTGenerator{
		Random: util.NewRandom(),
	}
}

type ValueGenerator struct {
	*util.Random
}

func NewValueGenerator() *ValueGenerator {
	return &ValueGenerator{
		Random: util.NewRandom(),
	}
}

func (g *ValueGenerator) newRandomString(len int) (str string) {
	for i := 0; i < len; i++ {
		str = fmt.Sprintf("%s%s", str, string(
			rune(
				g.RandomRange(
					int([]rune("a")[0]),
					int([]rune("{")[0]),
				),
			),
		))
	}

	return str
}

func (g *ValueGenerator) NewString() string {
	switch g.RandomNum(64) {
	case 0:
		return strings.ToUpper(g.newRandomString(g.RandomNum(1024)))
	case 1:
		return strings.ToUpper(g.newRandomString(g.RandomNum(512)))
	case 2:
		return strings.ToUpper(g.newRandomString(g.RandomNum(256)))
	default:
		return strings.ToUpper(g.newRandomString(g.RandomNum(8)))
	}
}

func (g *ValueGenerator) NewInt() int64 {
	switch g.RandomNum(64) {
	case 0:
		return g.RandomInt63()
	case 1:
		return g.RandomInt63n(int64(2147483648))
	case 2:
		return int64(g.RandomNum(65536))
	default:
		return int64(g.RandomRange(0, 1024))
	}
}

func (g *ValueGenerator) NewFloat() float64 {
	return float64(g.NewInt()) * g.RandomFloat()
}

func (g *ValueGenerator) NewTimestamp() time.Time {
	min := time.Date(1970, 1, 1, 0, 0, 1, 0, time.UTC).Unix()
	max := time.Date(2100, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	return time.Unix(g.RandomInt63n(max-min)+min, 0)
}
