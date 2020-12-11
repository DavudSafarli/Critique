package contracts

import (
	"context"
	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/testing_utils"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"
)

// These are vars that contract implementations should override and run tests.
// This basically says that, if you want to implement me(contract) correctly:
//  - define struct that implements X Contract
//	- you should override X-ForTest variable in X_impl_test.go file
// 	- your tests should pass
var (
	StorageForTest = testcase.Var{
		Name: "StorageForTest",
		//Init: this variable should be initialized by each implementation using `StorageForTest.Let(spec, ...)`
	}
	GetStorage = func(t *testcase.T) Storage {
		return StorageForTest.Get(t).(Storage)
	}

	TxContextForTest = testcase.Var{
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

	GetTxContextForTest = func(t *testcase.T) context.Context {
		return TxContextForTest.Get(t).(context.Context)
	}
)

var Feedback = testcase.Var{
	Name: `feedbackVar`,
	Init: func(t *testcase.T) interface{} {
		f := testing_utils.ExampleFeedback()
		require.Nil(t, GetStorage(t).CreateFeedback(GetTxContextForTest(t), f))
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
