package impl

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/external/repository"
	"github.com/DavudSafarli/Critique/util/testing_utils"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"
)

func TestAttchUc(t *testing.T) {
	spec := testcase.NewSpec(t)
	driver, connstr := "pg", testing_utils.GetTestDbConnStr()
	attachmentRepo := repository.NewAttachmentRepository(driver, connstr)
	feedbackRepo := repository.NewFeedbackRepository(driver, connstr)
	attchuc := NewAttachmentUsecases(attachmentRepo).(AttachmentUsecasesImpl)
	feedbackuc := NewFeedbackUsecasesImpl(feedbackRepo, attachmentRepo)

	spec.After(func(t *testcase.T) {
		testing_utils.TruncateTestTables(t, "feedbacks", "attachments")
	})
	spec.Describe(`#CreateAttachments + #GetAttachments`, func(s *testcase.Spec) {
		input := testcase.Var{Name: "input"}
		feedbackID := testcase.Var{Name: "feedbackID"}
		var subjectCreator = func(t *testcase.T) ([]models.Attachment, error) {
			return attchuc.CreateAttachments(context.Background(), input.Get(t).([]models.Attachment), feedbackID.Get(t).(uint))
		}
		var subjectGetter = func(t *testcase.T) ([]models.Attachment, error) {
			fbid := feedbackID.Get(t).(uint)
			return attchuc.GetAttachments(context.Background(), fbid)
		}
		s.When(`input attachments are valid`, func(s *testcase.Spec) {
			input.Let(s, func(t *testcase.T) interface{} {
				return []models.Attachment{{Name: "smth", Path: "/smth"}}
			})
			// persist feedback beforehand, to make use of its ID
			s.Before(func(t *testcase.T) {
				validModel := models.Feedback{Title: "t", Body: "t", CreatedBy: "t", CreatedAt: uint(time.Now().Unix())}
				f, err := feedbackuc.CreateFeedback(context.Background(), validModel)
				require.Nil(t, err, "Create feedback should return no error")
				feedbackID.Set(t, f.ID)
			})
			s.Then(`values should persist`, func(t *testcase.T) {
				createdAttchs, err := subjectCreator(t)
				require.Nil(t, err)
				returnedAttchs, err := subjectGetter(t)
				require.Nil(t, err)
				inputLen := len(input.Get(t).([]models.Attachment))
				require.Equal(t, inputLen, len(returnedAttchs))
				reflect.DeepEqual(createdAttchs, returnedAttchs)
			})
		})
	})
}
