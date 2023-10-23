package data

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"testing"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = db.Debug()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestMigration(t *testing.T) {
	err := migrateDb(db)
	assert.NoError(t, err)
}
