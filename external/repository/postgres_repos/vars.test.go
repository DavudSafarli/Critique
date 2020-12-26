package postgres_repos

import (
	"sync"

	"github.com/DavudSafarli/Critique/domain/contracts"
	"github.com/DavudSafarli/Critique/testing_utils"
	"github.com/adamluzsi/testcase"
)

var (
	sharedGlobalDatabaseInstanceInit sync.Once
	sharedGlobalDatabaseInstance     database
	GetSharedGlobalDatabaseInstance  = func() database {
		sharedGlobalDatabaseInstanceInit.Do(func() {
			connStr := testing_utils.GetTestDbConnStr()
			storage, err := NewPostgresDatabase(connStr)
			if err != nil {
				panic("cannot initialize NewPostgresDatabase")
			}
			sharedGlobalDatabaseInstance = storage
		})
		return sharedGlobalDatabaseInstance
	}
)

func SetupPostgresStorageVar(s *testcase.Spec) {
	contracts.StorageForTest.Let(s, func(t *testcase.T) interface{} {
		return NewPostgresStorage(
			NewPGAttachmentRepository(GetSharedGlobalDatabaseInstance()),
			NewPGFeedbackRepository(GetSharedGlobalDatabaseInstance()),
			NewPGTagRepository(GetSharedGlobalDatabaseInstance()),
		)
	})
}
