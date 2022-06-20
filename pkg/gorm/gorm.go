package gorm

import (
	"context"
	"fmt"
	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	sgorm "github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
	gorm2 "github.com/goxiaoy/uow/gorm"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	gormKind = "gorm"
)

func NewConnStrResolver(c *conf.Endpoints, ts common.TenantStore) data.ConnStrResolver {
	conn := make(data.ConnStrings, 1)
	for k, v := range c.Databases {
		conn[k] = v.Source
	}
	mr := common.NewMultiTenancyConnStrResolver(func() common.TenantStore {
		return common.NewCachedTenantStore(ts)
	}, data.NewConnStrOption(conn))
	return mr
}

type DbProviderFunc func(ctx context.Context, key string) *gorm.DB

func (d DbProviderFunc) Get(ctx context.Context, key string) *gorm.DB {
	return d(ctx, key)
}

type ensureDbExistFunc func(string) error
type dbGuardianKey struct{}

// NewDbGuardianContext  flag for database auto creation
func NewDbGuardianContext(ctx context.Context, enable ...bool) context.Context {
	v := true
	if len(enable) > 0 {
		v = enable[0]
	}
	return context.WithValue(ctx, dbGuardianKey{}, v)
}

func isDbGuardianEnabled(ctx context.Context) bool {
	if common.GetMultiTenantSide(ctx) == common.Host {
		return true
	}
	if v, ok := ctx.Value(dbGuardianKey{}).(bool); ok {
		return v
	}
	return false
}

type DbCache struct {
	*common.Cache[string, *sgorm.DbWrap]
	d *conf.Data
}

// NewDbCache create a shared gorm.Db cache by dsn
func NewDbCache(d *conf.Data) (*DbCache, func()) {
	c := common.NewCache[string, *sgorm.DbWrap]()
	return &DbCache{Cache: c, d: d}, func() {
		c.Flush()
	}
}

func (c *DbCache) GetOrSet(ctx context.Context, key, connStr string) (*gorm.DB, error) {

	client, _, err := c.Cache.GetOrSet(connStr, func() (*sgorm.DbWrap, error) {
		//find config
		dbConfig := c.d.Endpoints.GetDatabaseMergedDefault(key)
		var dbGuardian ensureDbExistFunc
		//generate db
		tp := ""
		if dbConfig.TablePrefix == nil {
			tp = fmt.Sprintf("kit_%s_", key)
		} else {
			tp = dbConfig.TablePrefix.Value
		}
		var dial gorm.Dialector
		switch dbConfig.Driver {
		case "sqlite":
			dial = sqlite.Open(connStr)
		case "mysql":
			dial = mysql.Open(connStr)
			dbGuardian = func(s string) error {
				dsn, err := mysql2.ParseDSN(s)
				if err != nil {
					return err
				}
				dbname := dsn.DBName
				dsn.DBName = ""
				//open without db name
				db, err := gorm.Open(mysql.Open(dsn.FormatDSN()), &gorm.Config{})
				if err != nil {
					return err
				}
				if err := db.Use(otelgorm.NewPlugin()); err != nil {
					return err
				}
				err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbname)).Error
				if err != nil {
					return err
				}
				return closeDb(db)
			}
		default:
			panic("driver unsupported")
		}

		if isDbGuardianEnabled(ctx) && dbGuardian != nil {
			if err := dbGuardian(connStr); err != nil {
				return nil, err
			}
		}

		gormConf := &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: tp,
			}}
		client, err := gorm.Open(dial, gormConf)
		if err != nil {
			return nil, err
		}
		//register global
		RegisterCallbacks(client)
		if err := client.Use(otelgorm.NewPlugin(otelgorm.WithoutQueryVariables())); err != nil {
			panic(err)
		}
		if dbConfig.Debug {
			client = client.Debug()
		}
		return sgorm.NewDbWrap(client), nil
	})
	if err != nil {
		return nil, err
	}
	return client.WithContext(ctx), nil
}

func NewDbProvider(cache *DbCache, cs data.ConnStrResolver, d *conf.Data) sgorm.DbProvider {
	return DbProviderFunc(func(ctx context.Context, key string) *gorm.DB {
		//find connection string
		s, err := cs.Resolve(ctx, key)
		if err != nil {
			panic(err)
		}

		//find transactional db from uow
		if u, ok := uow.FromCurrentUow(ctx); ok {
			tx, err := u.GetTxDb(ctx, []string{gormKind, key, s})
			if err != nil {
				panic(err)
			}
			g, ok := tx.(*gorm2.TransactionDb)
			if !ok {
				panic(fmt.Errorf("%s is not a *gorm.DB instance", s))
			}
			return g.WithContext(ctx)
		}

		client, err := cache.GetOrSet(ctx, key, s)

		if err != nil {
			panic(err)
		}
		return client.WithContext(ctx)
	})

}

func closeDb(d *gorm.DB) error {
	sqlDB, err := d.DB()
	if err != nil {
		return err
	}
	cErr := sqlDB.Close()
	if cErr != nil {
		//todo logging
		//logger.Errorf("Gorm db close error: %s", err.Error())
		return cErr
	}
	return nil
}
