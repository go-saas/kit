package data

import (
	"context"
	"embed"
	"fmt"
	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"
	"github.com/go-saas/kit/pkg/dal"
	kitgorm "github.com/go-saas/kit/pkg/gorm"
	"github.com/go-saas/saas/seed"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
)

const (
	barrier = "dtmcli.barrier"
	storage = "dtmsvr.storage"
)

type migrator struct {
	provider dal.ConstDbProvider
	connName dal.ConnName
	kind     []string
}

type (
	BarrierMigrator struct {
		*migrator
	}
	StorageMigrator struct {
		*migrator
	}
	Migrator struct {
		*migrator
	}
)

//go:embed  sqls
var sqls embed.FS

func NewBarrierMigrator(provider dal.ConstDbProvider, connName dal.ConnName) *BarrierMigrator {
	return &BarrierMigrator{newMigrator(provider, connName, barrier)}
}
func NewStorageMigrator(provider dal.ConstDbProvider, connName dal.ConnName) *StorageMigrator {
	return &StorageMigrator{newMigrator(provider, connName, storage)}
}

func NewMigrator(provider dal.ConstDbProvider, connName dal.ConnName) *Migrator {
	return &Migrator{newMigrator(provider, connName, storage, barrier)}
}

func newMigrator(provider dal.ConstDbProvider, connName dal.ConnName, kind ...string) *migrator {
	return &migrator{
		provider: provider,
		connName: connName,
		kind:     kind,
	}
}

var _ seed.Contrib = (*migrator)(nil)

func (m *migrator) Seed(ctx context.Context, sCtx *seed.Context) error {
	skipDrop := true
	if len(sCtx.TenantId) == 0 {
		ctx = kitgorm.NewDbGuardianContext(ctx)
		//only apply for host
		//get db kind
		db := GetDb(ctx, m.provider, m.connName)
		name := db.Dialector.Name()
		dtmimp.SetCurrentDBType(name)
		//read sql
		fileNames := make([]string, len(m.kind))
		for i, k := range m.kind {
			fileNames[i] = fmt.Sprintf("sqls/%s.%s.sql", k, name)
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
func (m *migrator) do(ctx context.Context, fileName string, skipDrop bool, db *gorm.DB) error {
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
