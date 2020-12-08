package spec_helper

import (
	"context"
	"sync"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/domain/usecases/attachment_usecases"
	"github.com/DavudSafarli/Critique/domain/usecases/feedback_usecases"
	"github.com/DavudSafarli/Critique/external/repository/abstract"
	"github.com/DavudSafarli/Critique/external/repository/postgres_repos"
	"github.com/DavudSafarli/Critique/impl"
	"github.com/DavudSafarli/Critique/testing_utils"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"
)

// database/databases
var (
	sharedGlobalStorageInstanceInit sync.Once
	sharedGlobalStorageInstance     *postgres_repos.Storage

	getSharedGlobalStorageInstance = func() *postgres_repos.Storage {
		sharedGlobalStorageInstanceInit.Do(func() {
			connStr := testing_utils.GetTestDbConnStr()
			storage, err := postgres_repos.NewPostgresStorage(connStr)
			if err != nil {
				panic("cannot initialize NewPostgresStorage")
			}
			sharedGlobalStorageInstance = storage
		})
		return sharedGlobalStorageInstance
	}
	StorageForTest = testcase.Var{
		Name: `StorageForTest`,
		Init: func(t *testcase.T) interface{} {
			return getSharedGlobalStorageInstance()
		},
	}
	GetStorage = func(t *testcase.T) *postgres_repos.Storage {
		return StorageForTest.Get(t).(*postgres_repos.Storage)
	}
)

// initializations for different repository implementations
var (
	ContextForTx = testcase.Var{
		Name: "Context",
		Init: func(t *testcase.T) interface{} {
			ctx, err := GetStorage(t).BeginTx(context.Background())
			if err != nil {
				panic(`failed to begin Transaction`)
			}
			t.Defer(func() {
				GetStorage(t).RollbackTx(ctx)
			})
			return ctx
		},
	}

	TagRepoForTest = testcase.Var{
		Name: "TagRepoForTest",
		Init: func(t *testcase.T) interface{} {
			return postgres_repos.NewPGTagRepository(GetStorage(t))
		},
	}
	AttachmentRepoForTest = testcase.Var{
		Name: "AttachmentRepoForTest",
		Init: func(t *testcase.T) interface{} {
			return postgres_repos.NewPGAttachmentRepository(GetStorage(t))
		},
	}
	FeedbackRepoForTest = testcase.Var{
		Name: "FeedbackRepoForTest",
		Init: func(t *testcase.T) interface{} {
			return postgres_repos.NewPGFeedbackRepository(GetStorage(t))
		},
	}
)

// there could me inmemory repo initializations
//var ()

// usecases
var (
	AttachmentUsecaseForTest = testcase.Var{
		Name: "AttachmentUsecaseForTest",
		Init: func(t *testcase.T) interface{} {
			return impl.NewAttachmentUsecases(GetAttachmentRepoForTest(t))
		},
	}
	FeedbackUsecaseForTest = testcase.Var{
		Name: "FeedbackUsecaseForTest",
		Init: func(t *testcase.T) interface{} {
			return impl.NewFeedbackUsecasesImpl(GetFeedbackRepoForTest(t), GetAttachmentRepoForTest(t))
		},
	}
)

// getters for repo
var (
	GetTxContextForTest = func(t *testcase.T) context.Context {
		return ContextForTx.Get(t).(context.Context)
	}
	GetTagRepoForTest = func(t *testcase.T) abstract.TagRepository {
		return TagRepoForTest.Get(t).(abstract.TagRepository)
	}
	GetAttachmentRepoForTest = func(t *testcase.T) abstract.AttachmentRepository {
		return AttachmentRepoForTest.Get(t).(abstract.AttachmentRepository)
	}
	GetFeedbackRepoForTest = func(t *testcase.T) abstract.FeedbackRepository {
		return FeedbackRepoForTest.Get(t).(abstract.FeedbackRepository)
	}
)

// getters for usecases
var (
	GetAttachmentUsecaseForTest = func(t *testcase.T) attachment_usecases.AttachmentUsecases {
		return AttachmentUsecaseForTest.Get(t).(attachment_usecases.AttachmentUsecases)
	}
	GetFeedbackUsecaseForTest = func(t *testcase.T) feedback_usecases.FeedbackUsecases {
		return FeedbackUsecaseForTest.Get(t).(feedback_usecases.FeedbackUsecases)
	}
)

var Feedback = testcase.Var{
	Name: `feedbackVar`,
	Init: func(t *testcase.T) interface{} {
		f := testing_utils.ExampleFeedback()
		require.Nil(t, GetFeedbackRepoForTest(t).Create(GetTxContextForTest(t), f))
		return f
	},
}
var GetFeedback = func(t *testcase.T) *models.Feedback {
	return Feedback.Get(t).(*models.Feedback)
}
var FeedbackID = testcase.Var{
	Name: `feedbackIDVar`,
	Init: func(t *testcase.T) interface{} {
		return Feedback.Get(t).(*models.Feedback).ID
	},
}
var GetFeedbackID = func(t *testcase.T) uint {
	return FeedbackID.Get(t).(uint)
}
