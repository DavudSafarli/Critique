package api

import (
	"context"
	"encoding/json"
	"github.com/DavudSafarli/Critique/domain/contracts"
	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/impl"
	"github.com/adamluzsi/testcase"
	"github.com/adamluzsi/testcase/httpspec"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var FeedbackCtrlForTest = testcase.Var{
	Name: `FeedbackCtrlForTest`,
	Init: func(t *testcase.T) interface{} {
		ctrl := NewFeedbackCtrl(impl.GetFeedbackUsecaseForTest(t))
		return ctrl
	},
}
var GetFeedbackCtrlForTest = func(t *testcase.T) *FeedbackCtrl {
	return FeedbackCtrlForTest.Get(t).(*FeedbackCtrl)
}

type server struct {
	dep func(http.ResponseWriter, *http.Request)
}
func(s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.dep(w, r)
}
func funcToHandler(dep func(http.ResponseWriter, *http.Request)) server {
	return server{dep:dep}
}

func TestFeedbackCtrl(t *testing.T) {
	s := testcase.NewSpec(t)
	s.Parallel()
	t.Parallel()
	impl.SetupUsecaseDependencies(s)
	httpspec.HandlerSpec(s, func(t *testcase.T) http.Handler {
		return GetHandler(*GetFeedbackCtrlForTest(t))
	})
	s.Before(func(t *testcase.T) {
		httpspec.Header(t).Set(`Content-Type`, `application/json`)
	})
	httpspec.LetContext(s, func(t *testcase.T) context.Context {
		return contracts.GetTxContextForTest(t)
	})

	s.Describe(`POST /feedbacks - create Attachment`, func(s *testcase.Spec) {
		httpspec.LetMethodValue(s, http.MethodPost)
		httpspec.LetPathValue(s, `/feedbacks`)

		var onSuccess = func(t *testcase.T) models.Feedback {
			rr := httpspec.ServeHTTP(t)
			require.Equal(t, http.StatusOK, rr.Code)
			var resp models.Feedback
			require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			return resp
		}

		s.When(`valid feedback is provided`, func(s *testcase.Spec) {
			httpspec.LetBody(s, func(t *testcase.T) interface{} {
				return createFeedbackRequest{Feedback: models.Feedback{
					Title:       "Example",
					Body:        "Body",
					CreatedBy:   "id of someone",
				}}
			})
			s.Then(`it should return 200 OK`, func(t *testcase.T) {
				response := onSuccess(t)
				require.NotZero(t, response.ID)
			})
		})
	})
}