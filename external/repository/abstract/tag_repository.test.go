package abstract

import (
	"context"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"

	"github.com/DavudSafarli/Critique/domain/models"
)

// TestTagRepositoryBehaviour does what its name says
func TestTagRepositoryBehaviour(t *testing.T, abstractRepo TagRepository, cleanupFunc func()) {
	t.Run("can CreateMany tags at once and Get them all", func(t *testing.T) {
		t.Cleanup(cleanupFunc)
		insertedTags := []models.Tag{
			{Name: "NewTagName1"},
			{Name: "NewTagName2"},
			{Name: "NewTagName3"},
		}
		numOfInsertedRows := len(insertedTags)
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		insertedTags, err := abstractRepo.CreateMany(ctx, insertedTags)
		require.Nil(t, err, "Failed to add Tags")

		selectedTags, err := abstractRepo.Get(ctx)
		require.Nil(t, err, "Failed to select all Tags")

		require.Equal(t, numOfInsertedRows, len(selectedTags), "Get didn't return all newly inserted Tags")

		require.True(t, reflect.DeepEqual(insertedTags, selectedTags), "Inserted Tags are not the same as selected ones database")
	})

	t.Run("can CreateMany tags and RemoveMany at once", func(t *testing.T) {
		t.Cleanup(cleanupFunc)
		// arrange
		insertedTags := []models.Tag{
			{Name: "NewTagName1"},
			{Name: "NewTagName2"},
			{Name: "NewTagName3"},
		}
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		insertedTags, _ = abstractRepo.CreateMany(ctx, insertedTags)

		// act
		tagIDs := make([]uint, 0, len(insertedTags))
		for _, val := range insertedTags {
			tagIDs = append(tagIDs, val.ID)
		}
		err := abstractRepo.RemoveMany(ctx, tagIDs)

		// assert
		require.Nil(t, err, "failed to perform RemoveMany")
		selectedTags, err := abstractRepo.Get(ctx)
		require.Nil(t, err, "failed to perform Get")
		require.Zero(t, len(selectedTags), "Did not remove all Tags. Database still contains some")
	})
}
