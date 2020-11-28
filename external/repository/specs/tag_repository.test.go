package specs

import (
	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/spec_helper"
	"github.com/DavudSafarli/Critique/testing_utils"
	"github.com/adamluzsi/testcase"
	"github.com/stretchr/testify/require"
)

// TestTagRepositoryBehaviour does what its name says
func TestTagRepositoryBehaviour(s *testcase.Spec) {
	getCtx := spec_helper.GetTxContextForTest

	s.Describe(`TagRepository#CreateMany + #Get`, func(s *testcase.Spec) {
		subject := func(t *testcase.T) ([]models.Tag, error) {
			tags := testing_utils.ExampleTagSlice(7)
			return tags, spec_helper.GetTagRepoForTest(t).CreateMany(getCtx(t), tags)
		}
		s.When(`Life is beautiful`, func(s *testcase.Spec) {
			s.Then(`It should be able to retrieve those tags`, func(t *testcase.T) {
				createdTags, err := subject(t)
				require.Nil(t, err)

				retrievedTags, err := spec_helper.GetTagRepoForTest(t).Get(getCtx(t))
				require.Nil(t, err)
				require.ElementsMatch(t, createdTags, retrievedTags)
			})
		})
	})
	s.Describe(`TagRepository#CreateMany + #RemoveMany`, func(s *testcase.Spec) {
		s.When(`RemoveMany works correctly`, func(s *testcase.Spec) {
			insertedTags := testing_utils.ExampleTagSlice(7)
			s.Before(func(t *testcase.T) {
				require.Nil(t, spec_helper.GetTagRepoForTest(t).CreateMany(getCtx(t), insertedTags))
			})
			s.Then(`It should be able to delete those tags`, func(t *testcase.T) {
				tagIDs := make([]uint, 0, len(insertedTags))
				for _, val := range insertedTags {
					tagIDs = append(tagIDs, val.ID)
				}
				require.Nil(t, spec_helper.GetTagRepoForTest(t).RemoveMany(getCtx(t), tagIDs))
				remeaningTags, err := spec_helper.GetTagRepoForTest(t).Get(getCtx(t))
				require.Nil(t, err)
				require.Len(t, remeaningTags, 0)
			})
		})
	})

}
