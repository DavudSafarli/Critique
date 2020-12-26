package impl

import (
	"github.com/DavudSafarli/Critique/domain/contracts"
	"github.com/DavudSafarli/Critique/external/repository/postgres_repos"
	"github.com/adamluzsi/testcase"
)

// setupUsecaseDependencies invokes the `variable setup function` of implementation of certain dependencies
// For ex: changing this to in_memory.SetupStorageVar(spec) would make `usecases tests` run with in memory storage implementation
func SetupUsecaseDependencies(spec *testcase.Spec) {
	postgres_repos.SetupPostgresStorageVar(spec)
}

var (
	AttachmentUsecaseForTest = testcase.Var{
		Name: "AttachmentUsecaseForTest",
		Init: func(t *testcase.T) interface{} {
			return NewAttachmentUsecases(contracts.GetStorage(t))
		},
	}
	GetAttachmentUsecaseForTest = func(t *testcase.T) AttachmentUsecasesImpl {
		return AttachmentUsecaseForTest.Get(t).(AttachmentUsecasesImpl)
	}
	FeedbackUsecaseForTest = testcase.Var{
		Name: "FeedbackUsecaseForTest",
		Init: func(t *testcase.T) interface{} {
			return NewFeedbackUsecasesImpl(contracts.GetStorage(t))
		},
	}
	GetFeedbackUsecaseForTest = func(t *testcase.T) FeedbackUsecasesImpl {
		return FeedbackUsecaseForTest.Get(t).(FeedbackUsecasesImpl)
	}
)
