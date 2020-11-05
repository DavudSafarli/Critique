package abstract

import (
	"context"
	"reflect"
	"testing"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/util"
)

// TestTagRepositoryBehaviour does what its name says
func TestTagRepositoryBehaviour(t *testing.T, abstractRepo TagRepository, cleanupFunc func() error) {
	t.Parallel()
	t.Run("CreateMany and Get", func(t *testing.T) {
		t.Cleanup(util.CreateCleanupWrapper(t, cleanupFunc))
		insertedTags := []models.Tag{
			{
				Name: "NewTagName1",
			},
			{
				Name: "NewTagName2",
			},
			{
				Name: "NewTagName3",
			},
		}
		numOfInsertedRows := len(insertedTags)
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)

		insertedTags, err := abstractRepo.CreateMany(ctx, insertedTags)
		if err != nil {
			t.Fatalf("TagRepository.CreateMany: Failed to add Tags: %s", err)
		}
		selectedTags, err := abstractRepo.Get(ctx)
		if err != nil {
			t.Fatalf("TagRepository.Get: Failed to select all Tags: %s", err)
		}
		if len(selectedTags) != numOfInsertedRows {
			t.Fatalf("TagRepository.Get: Didn't return all newly inserted Tags: %s", err)
		}
		if reflect.DeepEqual(insertedTags, selectedTags) == false {
			t.Fatal("Inserted Tags are not the same as selected ones database")
		}
	})

	t.Run("RemoveMany", func(t *testing.T) {
		t.Cleanup(util.CreateCleanupWrapper(t, cleanupFunc))
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
		if err != nil {
			t.Fatalf("TagRepository.RemoveMany: Failed to remove Tags: %s", err)
		}
		selectedTags, _ := abstractRepo.Get(ctx)
		if len(selectedTags) != 0 {
			t.Fatalf("TagRepository.RemoveMany: Did not remove all Tags. Database still contains some: %s", err)
		}
	})
}