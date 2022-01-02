package session

import (
	"context"

	"github.com/pingcap/tidb/infoschema"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/planner/core"
	"github.com/pingcap/tidb/sessionctx"
	"github.com/pingcap/tidb/types"
)

type TiDBPlanBuilder struct {
	ctx     context.Context
	session sessionctx.Context
	schema  infoschema.InfoSchema
}

func NewTiDBPlanBuilder(tables []*model.TableInfo) *TiDBPlanBuilder {
	return &TiDBPlanBuilder{
		ctx:     context.Background(),
		session: core.MockContext(),
		schema:  infoschema.MockInfoSchema(tables),
	}
}

func (b *TiDBPlanBuilder) Build(stmt ast.Node) (core.Plan, types.NameSlice, error) {
	return core.BuildLogicalPlanForTest(b.ctx, b.session, stmt, b.schema)
}
