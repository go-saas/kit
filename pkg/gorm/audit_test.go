package gorm

import (
	"context"
	"github.com/go-saas/kit/pkg/authn"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type TestEntity struct {
	ID uint `gorm:"primaryKey,autoIncrement"`
	AuditedModel
}

func TestAudit(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?loc=Local"))
	assert.NoError(t, err)
	db = db.Debug()
	err = db.AutoMigrate(&TestEntity{})
	assert.NoError(t, err)

	RegisterAuditCallbacks(db)

	ctx := authn.NewUserContext(context.Background(), authn.NewUserInfo("test"))

	entity := &TestEntity{}

	err = db.WithContext(ctx).Create(entity).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, entity.ID)
	assert.NotEmpty(t, entity.CreatedAt)
	if assert.NotNil(t, entity.CreatedBy) {
		assert.Equal(t, "test", *entity.CreatedBy)
	}

	err = db.WithContext(ctx).Select("*").Updates(entity).Error
	assert.NoError(t, err)
	assert.NotEmpty(t, entity.UpdatedAt)
	if assert.NotNil(t, entity.UpdatedBy) {
		assert.Equal(t, "test", *entity.UpdatedBy)
	}

	dbEntity := &TestEntity{}
	now := time.Now()

	err = db.First(dbEntity, "created_at <= ?", &now).Error
	assert.NoError(t, err)
	assert.Equal(t, entity.ID, dbEntity.ID)
	if assert.NotNil(t, dbEntity.CreatedBy) {
		assert.Equal(t, "test", *dbEntity.CreatedBy)
	}
	if assert.NotNil(t, dbEntity.UpdatedBy) {
		assert.Equal(t, "test", *dbEntity.UpdatedBy)
	}
}
