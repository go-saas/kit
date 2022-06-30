package data

import (
	"context"
	"database/sql"
	"github.com/dtm-labs/dtmcli/dtmimp"
	"github.com/go-saas/kit/pkg/dal"
	"gorm.io/gorm"
)

func GetDb(ctx context.Context, provider dal.ConstDbProvider, connName dal.ConnName) *gorm.DB {
	return provider.Get(ctx, string(connName))
}

// ToSQLDB get the sql.DB
func ToSQLDB(db *gorm.DB) *sql.DB {
	d, err := db.DB()
	dtmimp.E2P(err)
	return d
}
