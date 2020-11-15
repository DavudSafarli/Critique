package testing_utils

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Get connection string for LOCAL TEST DATABASE
func GetTestDbConnStr() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("local-test-user", "local-test-pass"),
		Host:     "localhost:5433",
		Path:     "critique-local-test",
		RawQuery: "sslmode=disable",
	}
	connStr := u.String()
	return connStr
}
func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PrepareTestDb() {
	fmt.Println("PrepareTestDb PrepareTestDb  PrepareTestDb PrepareTestDb PrepareTestDb PrepareTestDb PrepareTestDb PrepareTestDb PrepareTestDb PrepareTestDb PrepareTestDb PrepareTestDb PrepareTestDb PrepareTestDb")
	connStr := GetTestDbConnStr()
	migrationsFolder := "file://D:/GO/GOPATH/go/src/github.com/DavudSafarli/Critique/scripts/db/migrations"
	m, err := migrate.New(migrationsFolder, connStr)
	panicOnErr(err)
	// Migrate all the way up ...
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}

func TruncateTestTables(t *testing.T, tableNamesSlice ...string) func() {
	return func() {
		tableNamesStr := strings.Join(tableNamesSlice, ", ")
		conn, err := pgx.Connect(context.Background(), GetTestDbConnStr())
		panicOnErr(err)
		sql := "TRUNCATE TABLE " + tableNamesStr + " RESTART IDENTITY CASCADE"
		_, err = conn.Exec(context.Background(), sql)
		panicOnErr(err)
		conn.Close(context.Background())
	}
}

func CreateCleanupWrapper(t *testing.T, cleanupFunc func() error) func() {
	return func() {
		err := cleanupFunc()
		if err != nil {
			t.Errorf("error in cleanup function for %s: %s", t.Name(), err)
		}
	}
}
