package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func GetHandler(feedbackCtrl FeedbackCtrl) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/feedbacks", feedbackCtrl.Create)
	r.Get("/feedbacks/{id}", feedbackCtrl.Get)
	return r
}
