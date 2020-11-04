package main

import (
	"context"
	"time"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/external/repository"
	"github.com/DavudSafarli/Critique/impl"
)

// main starts the server
func main() {
	// get it from config or through env variables
	connstr := "postgres://admin:critiquesecretpassword@localhost/critique?sslmode=disable"

	feedbackrepo := repository.NewFeedbackRepository("pg", connstr)
	feedbackvalidator := models.FeedbackValidator{}
	feedbackImpls := impl.NewFeedbackUsecasesImpl(feedbackrepo, feedbackvalidator)
	// can be used and passed to anywhere seperately
	//creator := feedbackImpls.(feedback_usecases.FeedbackCreator)
	//getter := feedbackImpls.(feedback_usecases.FeedbackDetailsGetter)
	//paginator := feedbackImpls.(feedback_usecases.FeedbackPaginator)

	// create new Feedback
	f := models.Feedback{
		Body:      "a",
		CreatedBy: "",
		CreatedAt: uint(time.Now().Unix()),
	}
	// use usecase
	feedbackImpls.CreateFeedback(context.Background(), f)

}
