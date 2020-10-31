package repositories

import (
	"testing"

	"github.com/DavudSafarli/Critique/pkg/api/domaintest"
	"github.com/DavudSafarli/Critique/pkg/database/postgres"
	"github.com/DavudSafarli/Critique/pkg/util"
)

func CreateStorage(t *testing.T) *postgres.Storage {
	pgConnectionString := util.RunPostgresDockerAndGetConnectionString(t)
	storage, err := postgres.NewDbConnection(pgConnectionString)
	util.MigrateDatabase(t, pgConnectionString)
	if err != nil {
		t.Fatalf("Failed to create a new Storage: %s", err)
	}
	return storage
}

func TestRepository(t *testing.T) {
	storage := CreateStorage(t)
	repo := &TagRepository{
		storage: storage,
	}
	domaintest.TestTagRepositoryBehaviour(t, repo)
}
