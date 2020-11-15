package postgres_repos

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//
//func TestMain(m *testing.M) {
//	testing_utils.PrepareTestDb()
//	code := m.Run()
//	os.Exit(code)
//}
