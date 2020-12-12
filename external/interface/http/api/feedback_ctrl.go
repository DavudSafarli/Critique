package api

import (
	"encoding/json"
	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/domain/usecases/feedback_usecases"
	"net/http"
)

type FeedbackCtrl struct {
	usecases feedback_usecases.FeedbackUsecases
}
func NewFeedbackCtrl(usecases feedback_usecases.FeedbackUsecases) *FeedbackCtrl{
	return &FeedbackCtrl{
		usecases: usecases,
	}
}

type createFeedbackRequest struct {
	Feedback models.Feedback
}

func (ctrl *FeedbackCtrl) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var req createFeedbackRequest
	decoder.Decode(&req)

	ctx := r.Context()
	feedback, err := ctrl.usecases.CreateFeedback(ctx, req.Feedback)
	if err != nil {

	}
	responseJson(w, feedback)
}

func responseJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}