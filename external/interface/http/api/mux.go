package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func GetHandler(feedbackCtrl FeedbackCtrl) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/feedbacks", feedbackCtrl.Create)
	return r
}