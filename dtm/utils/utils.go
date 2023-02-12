package utils

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	gorm2 "github.com/go-saas/uow/gorm"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v4"
	"gorm.io/gorm"
	"net"
	"strconv"
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

func ParseDsnToDbConfig(driver, dsn string) (*dtmcli.DBConf, error) {
	ret := &dtmcli.DBConf{
		Driver: driver,
	}
	switch driver {
	case "mysql":
		cfg, err := mysql.ParseDSN(dsn)
		if err != nil {
			return nil, err
		}
		ret.User = cfg.User
		ret.Password = cfg.Passwd
		ret.Db = cfg.DBName
		host, port, err := net.SplitHostPort(cfg.Addr)
		if len(host) == 0 {
			host = "localhost"
		}
		if len(port) == 0 {
			port = "3306"
		}
		if err != nil {
			return nil, err
		}
		ret.Host = host
		if len(port) > 0 {
			port, _ := strconv.ParseInt(port, 10, 64)
			ret.Port = port
		}
	case "postgres":
		cfg, err := pgx.ParseConfig(dsn)
		if err != nil {
			return nil, err
		}
		ret.User = cfg.User
		ret.Password = cfg.Password
		ret.Db = cfg.Database
		ret.Host = cfg.Host
		ret.Port = int64(cfg.Port)
	default:
		panic("driver unsupported")
	}

	return ret, nil

}
