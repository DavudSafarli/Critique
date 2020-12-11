package contracts

import (
	"context"
	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/testing_utils"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"
	"math"
)

type FeedbackRepoitoryTester interface {
	FeedbackRepository
	GetAll(ctx context.Context) ([]models.Feedback, error)
}

// TestAttachmentRepositoryBehaviour does what its name says
func TestFeedbackRepositoryBehaviour(s *testcase.Spec) {
	getCtx := GetTxContextForTest

	s.Describe(`FeedbackRepository#Create`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) models.Feedback {
			feedback := testing_utils.ExampleFeedback()
			require.Nil(t, GetStorage(t).CreateFeedback(getCtx(t), feedback))
			return *feedback
		}
		s.When(`Life is good`, func(s *testcase.Spec) {
			s.Then(`It should be retrievable`, func(t *testcase.T) {
				inserteds := []models.Feedback{subject(t), subject(t), subject(t)}
				retrieved, err := GetStorage(t).(FeedbackRepoitoryTester).GetAll(getCtx(t))
				require.Nil(t, err)
				require.ElementsMatch(t, inserteds, retrieved)
			})
		})
		// I can't think any edge cases, because Repositories will receive valid models
	})

	s.Describe(`FeedbackRepository#Find`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) (models.Feedback, error) {
			return GetStorage(t).FindFeedback(getCtx(t), GetFeedbackID(t))
		}
		s.When(`There are multiple feedbacks beforehand`, func(s *testcase.Spec) {
			FeedbackID.Let(s, nil) //
			s.Before(func(t *testcase.T) {
				require.Nil(t, GetStorage(t).CreateFeedback(getCtx(t), testing_utils.ExampleFeedback()))
				require.Nil(t, GetStorage(t).CreateFeedback(getCtx(t), testing_utils.ExampleFeedback()))
			})
			s.Then(`It should the correct one`, func(t *testcase.T) {
				foundFeedback, err := subject(t)
				require.Nil(t, err)
				require.Equal(t, foundFeedback, *GetFeedback(t))
			})
		})
		// TODO: add case for feedback with attachment
	})

	s.Describe(`FeedbackRepository#GetFeedbacksPaginated`, func(s *testcase.Spec) {
		skip := s.LetValue(`skip`, uint(0))
		limit := s.LetValue(`limit`, uint(10))

		subject := func(t *testcase.T) ([]models.Feedback, error) {
			return GetStorage(t).GetFeedbacksPaginated(getCtx(t), skip.Get(t).(uint), limit.Get(t).(uint))
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
				require.Nil(t, GetStorage(t).CreateFeedback(getCtx(t), f))
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
					require.Len(t, selectedFeedbacks, lengthShouldBe, "GetFeedbacksPaginated did not return correct size slice")

					require.ElementsMatch(t, selectedFeedbacks, inserteds[skipInt:skipInt+lengthShouldBe])
					skipInt += 4 // increment randomly
				}
			})
		})
	})

}

type AttachmentRepositoryTester interface {
	AttachmentRepository
	GetAllAttachments(ctx context.Context) ([]models.Attachment, error)
}

// TestAttachmentRepositoryBehfeedbackIDaviour does what its name says
func TestAttachmentRepositoryBehaviour(s *testcase.Spec) {
	getCtx := GetTxContextForTest
	feedbackID := FeedbackID
	// creates  new feedback and returns its id
	s.Describe(`AttachmentRepository#CreateMany`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) ([]models.Attachment, error) {
			slice := testing_utils.ExampleAttchSlice(3)
			return slice, GetStorage(t).CreateManyAttachments(getCtx(t), slice, feedbackID.Get(t).(uint))
		}
		s.When(`There is feedback beforehand`, func(s *testcase.Spec) {
			feedbackID.Let(s, nil) // <-- bind to the scope
			s.Then(`It should not return error and be able to retrieve`, func(t *testcase.T) {
				inserteds, err := subject(t)
				require.Nil(t, err)
				retrieveds, err := GetStorage(t).(AttachmentRepositoryTester).GetAllAttachments(getCtx(t))
				require.Nil(t, err)
				require.Equal(t, inserteds, retrieveds)
			})
		})
		s.When(`There is no feedback`, func(s *testcase.Spec) {
			feedbackID.LetValue(s, uint(999))
			s.Then(`It should return error, because creating attachments is impossible`, func(t *testcase.T) {
				_, err := subject(t)
				require.Error(t, err)
			})
		})
	})

	s.Describe(`AttachmentRepository#GetByFeedbackID`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) ([]models.Attachment, error) {
			return GetStorage(t).GetAttachmentsByFeedbackID(getCtx(t), feedbackID.Get(t).(uint))
		}
		s.When(`There are multiple feedbacks and attachments beforehand`, func(s *testcase.Spec) {
			feedbackID.Let(s, nil) // <-- bind it to scope
			s.Before(func(t *testcase.T) {
				i := feedbackID.Get(t).(uint)
				require.Nil(t, GetStorage(t).CreateManyAttachments(getCtx(t), testing_utils.ExampleAttchSlice(2), i))

				feedback2 := testing_utils.ExampleFeedback()
				require.Nil(t, GetStorage(t).CreateFeedback(getCtx(t), feedback2))
				require.Nil(t, GetStorage(t).CreateManyAttachments(getCtx(t), testing_utils.ExampleAttchSlice(3), feedback2.ID))
			})
			s.Then(`It should find them`, func(t *testcase.T) {
				attchs, err := subject(t)
				require.Nil(t, err, "No error pls")
				require.Len(t, attchs, 2)
			})
		})
		s.When(`There is feedbacks and but no attachments`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				f := testing_utils.ExampleFeedback()
				require.Nil(t, GetStorage(t).CreateFeedback(getCtx(t), f))
				feedbackID.Set(t, f.ID)
			})
			s.Then(`should return empty slice without error`, func(t *testcase.T) {
				attchs, err := subject(t)
				require.Nil(t, err, "No error pls")
				require.Len(t, attchs, 0)
			})
		})
		s.When(`There is no feedback`, func(s *testcase.Spec) {
			s.Before(func(t *testcase.T) {
				feedbackID.Set(t, uint(0))
			})
			s.Then(`should return empty slice without error`, func(t *testcase.T) {
				attchs, err := subject(t)
				require.Nil(t, err)
				require.Len(t, attchs, 0)
			})
		})
	})
}

func TestTagRepositoryBehaviour(s *testcase.Spec) {
	getCtx := GetTxContextForTest

	s.Describe(`TagRepository#CreateMany + #Get`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) ([]models.Tag, error) {
			tags := testing_utils.ExampleTagSlice(7)
			return tags, GetStorage(t).CreateManyTags(getCtx(t), tags)
		}
		s.When(`Life is beautiful`, func(s *testcase.Spec) {
			s.Then(`It should be able to retrieve those tags`, func(t *testcase.T) {
				createdTags, err := subject(t)
				require.Nil(t, err)

				retrievedTags, err := GetStorage(t).GetTags(getCtx(t))
				require.Nil(t, err)
				require.ElementsMatch(t, createdTags, retrievedTags)
			})
		})
	})
	s.Describe(`TagRepository#CreateMany + #RemoveMany`, func(s *testcase.Spec) {
		s.When(`RemoveMany works correctly`, func(s *testcase.Spec) {
			insertedTags := testing_utils.ExampleTagSlice(7)
			s.Before(func(t *testcase.T) {
				require.Nil(t, GetStorage(t).CreateManyTags(getCtx(t), insertedTags))
			})
			s.Then(`It should be able to delete those tags`, func(t *testcase.T) {
				tagIDs := make([]uint, 0, len(insertedTags))
				for _, val := range insertedTags {
					tagIDs = append(tagIDs, val.ID)
				}
				require.Nil(t, GetStorage(t).RemoveManyTags(getCtx(t), tagIDs))
				remeaningTags, err := GetStorage(t).GetTags(getCtx(t))
				require.Nil(t, err)
				require.Len(t, remeaningTags, 0)
			})
		})
	})

}
