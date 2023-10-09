package gorm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-saas/kit/dtm/utils"
	"github.com/go-saas/kit/pkg/conf"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	sgorm "github.com/go-saas/saas/gorm"
	"github.com/go-saas/uow"
	gorm2 "github.com/go-saas/uow/gorm"
	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"reflect"
)

const (
	UowKind = "gorm"
)

func NewConnStrResolver(c *conf.Endpoints, ts saas.TenantStore) data.ConnStrResolver {
	conn := make(data.ConnStrings, 1)
	for k, v := range c.Databases {
		conn[k] = v.Source
	}
	mr := saas.NewMultiTenancyConnStrResolver(ts, conn)
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
	if v, ok := ctx.Value(dbGuardianKey{}).(bool); ok {
		return v
	}
	return false
}

// SqlDbCache adapter for dtm
type SqlDbCache struct {
	cache *saas.Cache[string, *sql.DB]
}

func NewSqlDbCache() (*SqlDbCache, func()) {
	c := saas.NewCache[string, *sql.DB]()
	return &SqlDbCache{cache: c}, func() {
		c.Flush()
	}
}

func (s *SqlDbCache) LoadOrStore(conf dtmimp.DBConf, factory func(conf dtmimp.DBConf) (*sql.DB, error)) (*sql.DB, error) {
	dsn := dtmimp.GetDsn(conf)
	db, _, err := s.cache.GetOrSet(dsn, func() (*sql.DB, error) {
		return factory(conf)
	})
	return db, err
}

type DbCache struct {
	cache *SqlDbCache
	d     *conf.Data
	l     klog.Logger
}

// NewDbCache create a shared gorm.Db cache by dsn
func NewDbCache(d *conf.Data, l klog.Logger, c *SqlDbCache) *DbCache {
	return &DbCache{cache: c, d: d, l: l}
}

func (c *DbCache) GetOrSet(ctx context.Context, key, dsn string) (*gorm.DB, error) {
	//find config
	dbConfig := c.d.Endpoints.GetDatabaseMergedDefault(key)
	dtmConf, err := utils.ParseDsnToDbConfig(dbConfig.Driver, dsn)
	if err != nil {
		return nil, err
	}
	sqlDb, err := c.cache.LoadOrStore(*dtmConf, func(conf dtmimp.DBConf) (*sql.DB, error) {
		var dbGuardian ensureDbExistFunc
		switch dbConfig.Driver {
		case "sqlite":
			//empty
		case "mysql":
			dbGuardian = func(s string) error {
				dsn, err := mysql2.ParseDSN(s)
				if err != nil {
					return err
				}
				dbname := dsn.DBName
				dsn.DBName = ""
				//open without db name
				db, err := gorm.Open(mysql.Open(dsn.FormatDSN()), GetDbConf(c.l, key, dbConfig))
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
			if err := dbGuardian(dsn); err != nil {
				return nil, err
			}
		}
		return dtmimp.StandaloneDB(conf)
	})
	if err != nil {
		return nil, err
	}
	client, err := OpenFromExisting(sqlDb, c.l, key, dbConfig)
	if err != nil {
		return nil, err
	}
	return client.WithContext(ctx), nil
}

func GetDbConf(l klog.Logger, key string, dbConfig *conf.Database) *gorm.Config {
	dbLogger := &Logger{
		Logger:   l,
		LogLevel: logger.Info,
	}
	return &gorm.Config{
		Logger: dbLogger}
}

func ApplyDefault(client *gorm.DB, dbConfig *conf.Database) *gorm.DB {
	ret := client
	//register global
	RegisterAuditCallbacks(client)
	RegisterAggCallbacks(client)
	if err := ret.Use(otelgorm.NewPlugin(otelgorm.WithoutQueryVariables())); err != nil {
		panic(err)
	}
	if dbConfig.Debug {
		ret = ret.Debug()
	}
	return ret
}

func OpenFromExisting(db *sql.DB, l klog.Logger, key string, dbConfig *conf.Database) (ret *gorm.DB, err error) {
	cfg := GetDbConf(l, key, dbConfig)
	switch dbConfig.Driver {
	case "sqlite":
		ret, err = gorm.Open(&sqlite.Dialector{Conn: db})
	case "mysql":
		ret, err = gorm.Open(mysql.New(mysql.Config{
			Conn: db,
		}), cfg)
		if err != nil {
			return nil, err
		}
	case "postgres":
		ret, err = gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), cfg)
		if err != nil {
			return nil, err
		}
	default:
		panic("driver unsupported")
	}
	if ret != nil {
		ret = ApplyDefault(ret, dbConfig)
	}
	return
}

func NewDbProvider(cache *DbCache, cs data.ConnStrResolver, d *conf.Data) sgorm.DbProvider {
	return DbProviderFunc(func(ctx context.Context, key string) *gorm.DB {
		//find from context. for dtm usage
		db, ok := fromContext(ctx, contextDbKey(key))
		if ok && db != nil {
			return db
		}

		//find connection string
		s, err := cs.Resolve(ctx, key)
		if err != nil {
			panic(err)
		}

		//find transactional db from uow
		if u, ok := uow.FromCurrentUow(ctx); ok {
			tx, err := u.GetTxDb(ctx, UowKind, key, s)
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
		klog.Errorf("gorm db close error: %s", err.Error())
		return cErr
	}
	return nil
}

func IsModel[T any](db *gorm.DB) (t T, is bool) {
	if db.Statement.Model != nil {
		t, is = db.Statement.Model.(T)
		if is {
			return
		}
	}
	if db.Statement.Schema == nil || db.Statement.Schema.ModelType == nil {
		return
	}
	_, is = reflect.New(db.Statement.Schema.ModelType).Interface().(T)
	return
}
