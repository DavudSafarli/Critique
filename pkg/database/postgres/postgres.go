package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v4"
)

// Storage is storage
type Storage struct {
	DB *pgxpool.Pool
	SB squirrel.StatementBuilderType
}

// NewDbConnection connects to the DB of passed connectionString and creates new pgxpool connection to application database
//
// Using docker, you can run :
//
// docker run --name critique-db
// -e POSTGRES_USER=admin
// -e POSTGRES_PASSWORD=critiquesecretpassword
// -e POSTGRES_DB=critique
// -p 5432:5432 -d postgres
//
// And use connection string below, to connect to postgres database:
//
// postgres://admin:critiquesecretpassword@localhost/critique?sslmode=disable
func NewDbConnection(connStr string) (*Storage, error) {
	poolConfig, err := pgxpool.ParseConfig(connStr)
	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Unable to create connection pool", "error", err)
	}

	return &Storage{
		DB: pool,
		SB: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}
