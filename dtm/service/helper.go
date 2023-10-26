package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/dtm-labs/dtm/client/workflow"
	klog "github.com/go-kratos/kratos/v2/log"
	dtmapi "github.com/go-saas/kit/dtm/api"
	"github.com/go-saas/kit/dtm/utils"
	sapi "github.com/go-saas/kit/pkg/api"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/kit/pkg/dal"
	"github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/kit/pkg/tracers"
	"github.com/go-saas/saas/data"
	"github.com/go-saas/uow"
	gorm2 "github.com/go-saas/uow/gorm"
	g "gorm.io/gorm"
)

type Helper struct {
	tokenMgr sapi.TokenManager
	d        *conf.Data
	cs       data.ConnStrResolver
	cache    *gorm.DbCache
	l        klog.Logger
	uowMgr   uow.Manager
	apiOpt   *sapi.Option
}

func NewHelper(tokenMgr sapi.TokenManager, d *conf.Data, cs data.ConnStrResolver, cache *gorm.DbCache, l klog.Logger, uowMgr uow.Manager, apiOpt *sapi.Option) *Helper {
	return &Helper{tokenMgr: tokenMgr, d: d, cs: cs, cache: cache, l: l, uowMgr: uowMgr, apiOpt: apiOpt}
}

type XaGrpcLocalFunc func(ctx context.Context, xa *dtmgrpc.XaGrpc) error

func (h *Helper) NewSagaGrpc(gid string, opts ...dtmgrpc.TransBaseOption) *dtmgrpc.SagaGrpc {
	return dtmgrpc.NewSagaGrpc(sapi.WithDiscovery(dtmapi.ServiceName), gid, opts...)
}

func (h *Helper) XaGlobalTransaction2(ctx context.Context, gid string, custom func(*dtmgrpc.XaGrpc), xaFunc dtmgrpc.XaGrpcGlobalFunc) error {
	custom1 := func(grpc *dtmgrpc.XaGrpc) {
		grpc.BranchHeaders = h.propagateHeader(ctx)
		if custom != nil {
			custom(grpc)
		}
	}
	return dtmgrpc.XaGlobalTransaction2(sapi.WithDiscovery(dtmapi.ServiceName), gid, custom1, xaFunc)
}

func (h *Helper) XaLocalTransaction(ctx context.Context, key dal.ConnName, xaFunc XaGrpcLocalFunc) error {
	cfg, err := h.GetDbConfig(ctx, key)
	if err != nil {
		return err
	}
	return dtmgrpc.XaLocalTransaction(ctx, *cfg, func(db *sql.DB, xa *dtmgrpc.XaGrpc) error {
		ctx, err := h.WrapContext(ctx, key, db)
		if err != nil {
			return err
		}
		return xaFunc(ctx, xa)
	})
}

func (h *Helper) TccGlobalTransaction2(ctx context.Context, gid string, custom func(grpc *dtmgrpc.TccGrpc), tccFunc dtmgrpc.TccGlobalFunc) error {
	custom1 := func(grpc *dtmgrpc.TccGrpc) {
		grpc.BranchHeaders = h.propagateHeader(ctx)
		if custom != nil {
			custom(grpc)
		}
	}
	return dtmgrpc.TccGlobalTransaction2(sapi.WithDiscovery(dtmapi.ServiceName), gid, custom1, tccFunc)
}

func (h *Helper) NewMsgGrpc(ctx context.Context, gid string, opts ...dtmgrpc.TransBaseOption) *dtmgrpc.MsgGrpc {
	//add token into header to pass TrustedContext validator
	opts = append(opts, dtmgrpc.WithBranchHeaders(h.propagateHeader(ctx)))
	return dtmgrpc.NewMsgGrpc(sapi.WithDiscovery(dtmapi.ServiceName), gid, opts...)
}

func (h *Helper) WorkflowRegister(name string, handler workflow.WfFunc, custom ...func(wf *workflow.Workflow)) error {
	custom1 := func(wf *workflow.Workflow) {
		wf.BranchHeaders = h.propagateHeader(wf.Context)
	}
	custom = append(custom, custom1)
	return workflow.Register(name, handler, custom...)
}

func (h *Helper) WorkflowRegister2(name string, handler workflow.WfFunc2, custom ...func(wf *workflow.Workflow)) error {
	custom1 := func(wf *workflow.Workflow) {
		wf.BranchHeaders = h.propagateHeader(wf.Context)
	}
	custom = append(custom, custom1)
	return workflow.Register2(name, handler, custom...)
}

type BarrierDbFunc func(ctx context.Context, u *uow.UnitOfWork) (*sql.Tx, error)

// BarrierUow see dtmcli.BranchBarrier.Call
func (h *Helper) BarrierUow(ctx context.Context, bb *dtmcli.BranchBarrier, key dal.ConnName, fn func(ctx context.Context) error, opt ...*sql.TxOptions) (rerr error) {
	u, err := h.uowMgr.CreateNew(ctx, opt...)
	if err != nil {
		return err
	}
	//push into context
	ctx = uow.NewCurrentUow(ctx, u)
	//already transactional ,barrier db is managed by uow now
	//find connection string
	_, dsn, err := h.ResolveDsn(ctx, key)
	txn, err := u.GetTxDb(ctx, gorm.UowKind, string(key), dsn)
	if err != nil {
		return err
	}
	barrierDb := utils.ToSQLDB(txn.(*gorm2.TransactionDb).DB)

	//dtmcli.BranchBarrier.newBarrierID()
	bb.BarrierID++
	bid := fmt.Sprintf("%02d", bb.BarrierID)

	defer dtmimp.DeferDo(&rerr, func() error {
		return u.Commit()
	}, func() error {
		return u.Rollback()
	})

	originOp := map[string]string{
		dtmimp.OpCancel:     dtmimp.OpTry,    // tcc
		dtmimp.OpCompensate: dtmimp.OpAction, // saga
		dtmimp.OpRollback:   dtmimp.OpAction, // workflow
	}[bb.Op]

	originAffected, oerr := dtmimp.InsertBarrier(barrierDb, bb.TransType, bb.Gid, bb.BranchID, originOp, bid, bb.Op, bb.DBType, bb.BarrierTableName)
	currentAffected, rerr := dtmimp.InsertBarrier(barrierDb, bb.TransType, bb.Gid, bb.BranchID, bb.Op, bid, bb.Op, bb.DBType, bb.BarrierTableName)
	klog.Debugf("originAffected: %d currentAffected: %d", originAffected, currentAffected)

	if rerr == nil && bb.Op == dtmimp.MsgDoOp && currentAffected == 0 { // for msg's DoAndSubmit, repeated insert should be rejected.
		return dtmcli.ErrDuplicated
	}

	if rerr == nil {
		rerr = oerr
	}

	if (bb.Op == dtmimp.OpCancel || bb.Op == dtmimp.OpCompensate || bb.Op == dtmimp.OpRollback) && originAffected > 0 || // null compensate
		currentAffected == 0 { // repeated request or dangled request
		return
	}
	if rerr == nil {
		rerr = fn(ctx)
	}
	return
}

func (h *Helper) GetDbConfig(ctx context.Context, key dal.ConnName) (*dtmcli.DBConf, error) {
	//find connection string
	cfg, dsn, err := h.ResolveDsn(ctx, key)
	if err != nil {
		return nil, err
	}
	return utils.ParseDsnToDbConfig(cfg.Driver, dsn)
}

func (h *Helper) ResolveDsn(ctx context.Context, key dal.ConnName) (cfg *conf.Database, dsn string, err error) {
	//find connection string
	s, err := h.cs.Resolve(ctx, string(key))
	if err != nil {
		return nil, "", err
	}
	//find config
	dbConfig := h.d.Endpoints.GetDatabaseMergedDefault(string(key))
	return dbConfig, s, nil
}

func (h *Helper) propagateHeader(ctx context.Context) map[string]string {
	headers := sapi.HeaderCarrier(dtmapi.MustAddBranchHeader(ctx, h.tokenMgr))
	//inject header
	for _, contributor := range h.apiOpt.Propagators {
		//do not handle error
		contributor.Inject(ctx, headers)
		tracers.DefaultPropagator.Inject(ctx, headers)
	}
	return headers
}

func (h *Helper) WrapContext(ctx context.Context, key dal.ConnName, db *sql.DB) (context.Context, error) {
	cfg, _, err := h.ResolveDsn(ctx, key)
	if err != nil {
		return ctx, err
	}
	gormDb, err := gorm.OpenFromExisting(db, h.l, string(key), cfg)
	gormDb = gormDb.Session(&g.Session{SkipDefaultTransaction: true})
	if err != nil {
		return ctx, err
	}
	ctx = gorm.NewContext(ctx, string(key), gormDb)
	gormDb = gormDb.WithContext(ctx)
	return ctx, nil
}
