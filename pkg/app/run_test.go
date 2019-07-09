package app_test

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/saromanov/migr/pkg/app"
	"github.com/stretchr/testify/assert"
)

const (
	basicPath       = "../../testdata/basic"
	migrCreateTable = "CREATE TABLE IF NOT EXISTS migr( id SERIAL, version int8 UNIQUE, changes varchar(128), hash varchar(128), applied bool, error_message varchar(128), failed bool, status varchar(16), created_at int8 )"
)

var (
	appTest *app.App
	db      *sql.DB
)

func init() {
	os.Setenv("MIGR_PATH", basicPath)
	appTest = app.New("postgres", "migr_test", "migr_test", "migr_test", "migr_test", 5432)
	dbTmp, err := createTestTable()
	if err != nil {
		fmt.Printf("unable to init db: %v", err)
	}
	db = dbTmp
}

// Create table for tests
func createTestTable() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=migr_test password=migr_test dbname=migr_test")
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %v", err)
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error to ping db")
	}

	return db, nil
}

func createMigrTable(db *sql.DB) error {
	_, err := db.Exec(migrCreateTable)
	if err != nil {
		return fmt.Errorf("unable to create migr table: %v", err)
	}

	return nil
}

func dropTable(db *sql.DB, name string) error {
	_, err := db.Exec(fmt.Sprintf("DROP TABLE %s", name))
	if err != nil {
		return fmt.Errorf("unable to drop migr table: %v", err)
	}

	return nil
}

func removeMigrDirs(pathDir string) error {
	files, err := ioutil.ReadDir(pathDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if !strings.Contains(f.Name(), "migr") {
			continue

		}
		os.RemoveAll(path.Join([]string{pathDir, f.Name()}...))
	}

	return nil
}

func writeFile(t *testing.T, path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0644)
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	err := createMigrTable(db)
	assert.NoError(t, err)
	defer func() {
		err := dropTable(db, "migr")
		assert.NoError(t, err)
		err = removeMigrDirs("../../testdata/basic")
		assert.NoError(t, err)
	}()

	err = appTest.Create(t.Name())
	assert.NoError(t, err)
}

func TestRun(t *testing.T) {
	err := createMigrTable(db)
	assert.NoError(t, err)
	defer func() {
		err := dropTable(db, "migr")
		assert.NoError(t, err)
		err = removeMigrDirs("../../testdata/basic")
		assert.NoError(t, err)
	}()
	err = appTest.Create(t.Name())
	assert.NoError(t, err)
	err = appTest.Run(basicPath)
	assert.NoError(t, err)

	versions, err := appTest.GetMigrationsInfo()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(versions))

	for _, v := range versions {
		assert.Equal(t, true, v.Applied)
	}
}

func TestDown(t *testing.T) {
	err := createMigrTable(db)
	assert.NoError(t, err)
	defer func() {
		err := dropTable(db, "migr")
		assert.NoError(t, err)
		err = removeMigrDirs("../../testdata/basic")
		assert.NoError(t, err)
	}()
	err = appTest.Create(t.Name())
	assert.NoError(t, err)
	err = appTest.Run(basicPath)
	assert.NoError(t, err)

	err = appTest.Down(basicPath)
	assert.NoError(t, err)
	versions, err := appTest.GetMigrationsInfo()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(versions))

	for _, v := range versions {
		assert.Equal(t, false, v.Applied)
	}
}
