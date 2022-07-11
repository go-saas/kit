package data

import (
	"context"
	"embed"
	"fmt"
	"github.com/dtm-labs/dtmcli/dtmimp"
	"github.com/go-saas/kit/pkg/dal"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/saas/seed"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
)

type Migrator struct {
	provider dal.ConstDbProvider
	connName dal.ConnName
}

//go:embed  sqls
var sqls embed.FS

func NewMigrator(provider dal.ConstDbProvider, connName dal.ConnName) *Migrator {
	return &Migrator{
		provider: provider,
		connName: connName,
	}
}

var _ seed.Contrib = (*Migrator)(nil)

func (m *Migrator) Seed(ctx context.Context, sCtx *seed.Context) error {
	skipDrop := true
	if len(sCtx.TenantId) == 0 {
		ctx = kitgorm.NewDbGuardianContext(ctx)
		//only apply for host
		//get db kind
		db := GetDb(ctx, m.provider, m.connName)
		name := db.Dialector.Name()
		dtmimp.SetCurrentDBType(name)
		//read sql
		fileNames := []string{
			fmt.Sprintf("sqls/dtmcli.barrier.%s.sql", name),
			fmt.Sprintf("sqls/dtmsvr.storage.%s.sql", name),
		}
		for _, fileName := range fileNames {
			err := m.do(ctx, fileName, skipDrop, db)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func (m *Migrator) do(ctx context.Context, fileName string, skipDrop bool, db *gorm.DB) error {
	fs, err := sqls.Open(fileName)
	defer fs.Close()
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(fs)
	if err != nil {
		return err
	}
	sqls := strings.Split(string(content), ";")
	for _, sql := range sqls {
		s := strings.TrimSpace(sql)
		if s == "" || (skipDrop && strings.Contains(s, "drop")) {
			continue
		}
		err := db.Exec(s).Error
		if err != nil {
			return err
		}
	}
	return nil
}
