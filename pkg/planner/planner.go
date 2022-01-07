package planner

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pingcap/tidb/domain"
	"github.com/pingcap/tidb/executor"
	"github.com/pingcap/tidb/infoschema"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/model"
	"github.com/pingcap/tidb/planner/core"
	"github.com/pingcap/tidb/sessionctx"
	"github.com/pingcap/tidb/store/mockstore"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util/mock"
)

func init() {
	_ = executor.GlobalMemoryUsageTracker
	_ = executor.GlobalDiskUsageTracker
}

type TiDBPlanBuilder struct {
	ctx     context.Context
	session sessionctx.Context
	schema  infoschema.InfoSchema
}

func sessionCtx() sessionctx.Context {
	ctx := mock.NewContext()
	// driver := mockstore.EmbedUnistoreDriver{}
	// kvstore, err := driver.Open("unistore:///tmp/mock/tikv")
	driver := mockstore.MockTiKVDriver{}
	kvstore, err := driver.Open("mocktikv:///tmp/mock/tikv")
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx.Store = kvstore

	ctx.GetSessionVars().CurrentDB = "test"
	ctx.GetSessionVars().SnapshotTS = uint64((time.Now().UnixNano() / int64(time.Millisecond)) << 18)

	do := &domain.Domain{}
	if err := do.CreateStatsHandle(ctx); err != nil {
		panic(fmt.Sprintf("create mock context panic: %+v", err))
	}
	domain.BindDomain(ctx, do)

	return ctx
}

func NewTiDBPlanBuilder(tables []*model.TableInfo) *TiDBPlanBuilder {
	return &TiDBPlanBuilder{
		ctx:     context.Background(),
		session: sessionCtx(),
		schema:  infoschema.MockInfoSchema(tables),
	}
}

func (b *TiDBPlanBuilder) Build(stmt ast.Node) (core.Plan, types.NameSlice, error) {
	return core.BuildLogicalPlanForTest(b.ctx, b.session, stmt, b.schema)
}
