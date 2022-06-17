package gorm

import (
	"fmt"
	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/goxiaoy/go-saas-kit/pkg/conf"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	sgorm "github.com/goxiaoy/go-saas/gorm"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDbOpener() (sgorm.DbOpener, func()) {
	return sgorm.NewDbOpener(func(db *gorm.DB) *gorm.DB {
		RegisterCallbacks(db)
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			panic(err)
		}
		return db
	})
}

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

func NewDbProvider(cs data.ConnStrResolver, c *sgorm.Config, opener sgorm.DbOpener) sgorm.DbProvider {
	return sgorm.NewDefaultDbProvider(cs, c, opener)
}

//NewGormConfig generate config from database endpoint by name.
func NewGormConfig(d *conf.Data, name string) *sgorm.Config {
	databases := d.Endpoints
	var c *conf.Database
	var ok bool
	c, ok = databases.Databases[name]
	if !ok {
		c, ok = databases.Databases[data.Default]
	}
	tp := ""
	if c.TablePrefix == nil {
		tp = fmt.Sprintf("kit_%s_", name)
	} else {
		tp = c.TablePrefix.Value
	}
	cfg := &sgorm.Config{
		Debug: c.Debug,
		Cfg: &gorm.Config{NamingStrategy: schema.NamingStrategy{
			TablePrefix: tp,
		}},
	}
	if c.Driver == "mysql" {
		cfg.Dialect = func(s string) gorm.Dialector {
			return mysql.Open(s)
		}
		//set database guardian function
		cfg.EnsureDbExist = func(cfg *sgorm.Config, s string) error {
			dsn, err := mysql2.ParseDSN(s)
			if err != nil {
				return err
			}
			dbname := dsn.DBName
			dsn.DBName = ""
			//open without db name
			db, err := gorm.Open(cfg.Dialect(dsn.FormatDSN()), &gorm.Config{})
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
	}
	if c.Driver == "sqlite" {
		cfg.Dialect = func(s string) gorm.Dialector {
			return sqlite.Open(s)
		}
		i := 1
		//https://github.com/go-gorm/gorm/issues/2875
		cfg.MaxOpenConn = &i
		cfg.MaxIdleConn = &i
	}
	//TODO support more database by underlying gorm driver?
	return cfg
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
