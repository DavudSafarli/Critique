package specs

import (
	"context"
	"math"

	"github.com/DavudSafarli/Critique/external/repository/abstract"
	"github.com/DavudSafarli/Critique/spec_helper"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/testing_utils"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"
)

type FeedbackRepoitoryTester interface {
	abstract.FeedbackRepository
	GetAll(ctx context.Context) ([]models.Feedback, error)
}

// TestAttachmentRepositoryBehaviour does what its name says
func TestFeedbackRepositoryBehaviour(s *testcase.Spec) {
	getCtx := spec_helper.GetTxContextForTest

	s.Describe(`FeedbackRepository#Create`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) models.Feedback {
			feedback := testing_utils.ExampleFeedback()
			require.Nil(t, spec_helper.GetFeedbackRepoForTest(t).Create(getCtx(t), feedback))
			return *feedback
		}
		s.When(`Life is good`, func(s *testcase.Spec) {
			s.Then(`It should be retrievable`, func(t *testcase.T) {
				inserteds := []models.Feedback{subject(t), subject(t), subject(t)}
				retrieved, err := spec_helper.GetFeedbackRepoForTest(t).(FeedbackRepoitoryTester).GetAll(getCtx(t))
				require.Nil(t, err)
				require.ElementsMatch(t, inserteds, retrieved)
			})
		})
		// I can't think any edge cases, because Repositories will receive valid models
	})

	s.Describe(`FeedbackRepository#Find`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) (models.Feedback, error) {
			return spec_helper.GetFeedbackRepoForTest(t).Find(getCtx(t), spec_helper.GetFeedbackID(t))
		}
		s.When(`There are multiple feedbacks beforehand`, func(s *testcase.Spec) {
			spec_helper.FeedbackID.Let(s, nil) //
			s.Before(func(t *testcase.T) {
				require.Nil(t, spec_helper.GetFeedbackRepoForTest(t).Create(getCtx(t), testing_utils.ExampleFeedback()))
				require.Nil(t, spec_helper.GetFeedbackRepoForTest(t).Create(getCtx(t), testing_utils.ExampleFeedback()))
			})
			s.Then(`It should the correct one`, func(t *testcase.T) {
				foundFeedback, err := subject(t)
				require.Nil(t, err)
				require.Equal(t, foundFeedback, *spec_helper.GetFeedback(t))
			})
		})
		// TODO: add case for feedback with attachment
	})

	s.Describe(`FeedbackRepository#GetPaginated`, func(s *testcase.Spec) {
		skip := s.LetValue(`skip`, uint(0))
		limit := s.LetValue(`limit`, uint(10))

		subject := func(t *testcase.T) ([]models.Feedback, error) {
			return spec_helper.GetFeedbackRepoForTest(t).GetPaginated(getCtx(t), skip.Get(t).(uint), limit.Get(t).(uint))
		}
		s.When(`There are no feedbacks beforehand`, func(s *testcase.Spec) {
			s.Then(`It should return nil slice, no error`, func(t *testcase.T) {
				feedbacks, err := subject(t)
				require.Nil(t, err)
				require.Len(t, feedbacks, 0)
			})
		})
		s.When(`There are 10 feedbacks beforehand`, func(s *testcase.Spec) {
			createAssertGetFeedback := func(t *testcase.T) models.Feedback {
				f := testing_utils.ExampleFeedback()
				require.Nil(t, spec_helper.GetFeedbackRepoForTest(t).Create(getCtx(t), f))
				return *f
			}
			inserteds := []models.Feedback{}
			s.Before(func(t *testcase.T) {
				for i := 0; i < 10; i++ {
					inserteds = append(inserteds, createAssertGetFeedback(t))
				}
			})
			// If 10 inserted, then feched by 3-3, it should get them by 3-3-3-1
			s.Then(`It should return them in correct order.`, func(t *testcase.T) {
				var skipInt int = 0
				var limitInt int = 3
				limit.Set(t, uint(limitInt))
				for skipInt < 10 {
					skip.Set(t, uint(skipInt))
					selectedFeedbacks, err := subject(t)
					require.Nil(t, err)
					lengthShouldBe := int(math.Min(float64(len(inserteds)-skipInt), float64(limitInt)))
					require.Len(t, selectedFeedbacks, lengthShouldBe, "GetPaginated did not return correct size slice")

					require.ElementsMatch(t, selectedFeedbacks, inserteds[skipInt:skipInt+lengthShouldBe])
					skipInt += 4 // increment randomly
				}
			})
		})
	})

}
