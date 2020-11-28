package testing_utils

import (
	"fmt"
	"time"

	"github.com/DavudSafarli/Critique/domain/models"
)

// ExampleFeedback return example feedback
func ExampleFeedback() *models.Feedback {
	return &models.Feedback{
		Title:     "T",
		Body:      "T",
		CreatedBy: "T",
		CreatedAt: uint(time.Now().Unix()),
	}
}
func ExampleInvalidFeedback() *models.Feedback {
	return &models.Feedback{}
}

// ExampleFeedback return example feedback slice
func ExampleFeedbackSlice(n int) []models.Feedback {
	slice := []models.Feedback{}
	for i := 0; i < n; i++ {
		slice = append(slice, *ExampleFeedback())
	}
	return slice
}

// ExampleAttch return example attachment
func ExampleAttch() models.Attachment {
	return models.Attachment{
		Name: "A1",
		Path: "P1",
	}
}

// ExampleAttchSlice return example attachment slice
func ExampleAttchSlice(n int) []models.Attachment {
	slice := []models.Attachment{}
	for i := 0; i < n; i++ {
		slice = append(slice, ExampleAttch())
	}
	return slice
}

var tagCounter int = 0

// ExampleTag return example attachment
func ExampleTag() models.Tag {
	tagCounter++
	return models.Tag{
		Name: fmt.Sprintf("Name %d", tagCounter),
	}
}

// ExampleTagSlice return example attachment slice
func ExampleTagSlice(n int) []models.Tag {
	slice := []models.Tag{}
	for i := 0; i < n; i++ {
		slice = append(slice, ExampleTag())
	}
	return slice
}
