package api

import (
	"context"
	"encoding/json"
	"fmt"
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

func TestFeedbackCtrl(t *testing.T) {
	s := testcase.NewSpec(t)
	s.Parallel()
	t.Parallel()
	getCtx := contracts.GetTxContextForTest
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

	s.Describe(`POST /feedbacks - create a Feedback`, func(s *testcase.Spec) {
		httpspec.LetMethodValue(s, http.MethodPost)
		httpspec.LetPathValue(s, `/feedbacks`)

		var onSuccess = func(t *testcase.T) createFeedbackResponse {
			rr := httpspec.ServeHTTP(t)
			require.Equal(t, http.StatusOK, rr.Code)
			var resp createFeedbackResponse
			require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			require.NotEqual(t, models.Feedback{}, resp.Feedback)
			return resp
		}
		s.When(`valid feedback is provided`, func(s *testcase.Spec) {
			httpspec.LetBody(s, func(t *testcase.T) interface{} {
				return createFeedbackRequest{Feedback: models.Feedback{
					Title:     "Example",
					Body:      "Body",
					CreatedBy: "id of someone",
				}}
			})
			s.Then(`it should return 200 OK with a non-zero ID assigned`, func(t *testcase.T) {
				response := onSuccess(t)
				require.NotZero(t, response.Feedback.ID)
			})
			s.Then(`it should be able to find that created feedback `, func(t *testcase.T) {
				response := onSuccess(t)
				f, err := impl.GetFeedbackUsecaseForTest(t).GetFeedbackDetails(getCtx(t), response.Feedback.ID)
				require.Nil(t, err)
				require.Equal(t, response.Feedback.ID, f.ID)
			})
		})
		s.When(`input is invalid`, func(s *testcase.Spec) {
			httpspec.LetBody(s, func(t *testcase.T) interface{} {
				return createFeedbackRequest{Feedback: models.Feedback{}}
			})
			s.Then(`it should return 422`, func(t *testcase.T) {
				require.Equal(t, http.StatusUnprocessableEntity, httpspec.ServeHTTP(t).Code)
			})
		})
	})

	s.Describe(`GET /feedbacks/:id`, func(s *testcase.Spec) {
		httpspec.LetMethodValue(s, http.MethodGet)
		idVar := testcase.Var{Name: `idVar`}
		httpspec.LetPath(s, func(t *testcase.T) string {
			return fmt.Sprintf("/feedbacks/%d", idVar.Get(t).(uint))
		})

		var onSuccess = func(t *testcase.T) getFeedbackResponse {
			rr := httpspec.ServeHTTP(t)
			require.Equal(t, http.StatusOK, rr.Code)
			var resp getFeedbackResponse
			require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			return resp
		}
		var onFail = func(t *testcase.T) (ErrorResponse, int) {
			rr := httpspec.ServeHTTP(t)
			require.NotEqual(t, http.StatusOK, rr.Code)
			var resp ErrorResponse
			require.Nil(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			return resp, rr.Code
		}
		s.When(`existing feedback is asked`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				idVar.Set(t, contracts.GetFeedbackID(t))
			})
			s.Then(`is should respond with right feedback`, func(t *testcase.T) {
				response := onSuccess(t)
				require.NotZero(t, response.Feedback.ID)
				require.Equal(t, contracts.GetFeedback(t).ID, response.Feedback.ID) // unnecessary line
				require.Equal(t, *contracts.GetFeedback(t), response.Feedback)
			})
		})
		s.When(`non-existing feedback is asked`, func(s *testcase.Spec) {
			idVar.LetValue(s, uint(1))
			s.Then(`is should return 404`, func(t *testcase.T) {
				response, code := onFail(t)
				require.Equal(t, http.StatusNotFound, code)
				require.Equal(t, http.StatusNotFound, response.Status)
			})
		})
	})
}
