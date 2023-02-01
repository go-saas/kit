package utils

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/client/dtmcli"
	"github.com/dtm-labs/client/dtmcli/dtmimp"
	"github.com/dtm-labs/client/dtmgrpc"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	gorm2 "github.com/go-saas/uow/gorm"
	"gorm.io/gorm"
)

const (
	BarrierInvalid = "BARRIER_INVALID"
)

func BarrierFromContext(ctx context.Context) (*dtmcli.BranchBarrier, error) {
	if t, ok := transport.FromServerContext(ctx); ok {
		if ht, ok := t.(*khttp.Transport); ok {
			ba, err := dtmcli.BarrierFromQuery(ht.Request().URL.Query())
			if err != nil {
				return nil, errors.BadRequest(BarrierInvalid, err.Error())
			}
			return ba, nil
		}
		if _, ok := t.(*grpc.Transport); ok {
			ba, err := dtmgrpc.BarrierFromGrpc(ctx)
			if err != nil {
				return nil, errors.BadRequest(BarrierInvalid, err.Error())
			}
			return ba, err
		}
	}
	panic("unsupported context to find barrier")
}

// ToSQLDB get the sql.DB
func ToSQLDB(db *gorm.DB) *sql.DB {
	d, err := db.DB()
	dtmimp.E2P(err)
	return d
}

func ToSqlTx(db *gorm2.TransactionDb) *sql.Tx {
	return db.Statement.ConnPool.(*sql.Tx)
}
