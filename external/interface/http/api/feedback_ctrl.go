package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DavudSafarli/Critique/domain/models"
	"github.com/DavudSafarli/Critique/domain/usecases/feedback_usecases"
	"github.com/go-chi/chi"
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
type createFeedbackResponse struct {
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
		responseJSONError(w, err)
		return
	}
	responseJSON(w, createFeedbackResponse{Feedback: feedback})
}

type getFeedbackResponse struct {
	Feedback models.Feedback
}

func (ctrl *FeedbackCtrl) Get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return
	}
	ctx := r.Context()
	feedback, err := ctrl.usecases.GetFeedbackDetails(ctx, uint(id))
	if err != nil {
		responseJSONError(w, err)
		return
	}
	responseJSON(w, getFeedbackResponse{Feedback: feedback})
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// responseJSONError converts error into meaningful response body
// and sends response with proper status code
func responseJSONError(w http.ResponseWriter, err error) {
	response := toHttp(err)
	w.WriteHeader(response.Status)
	responseJSON(w, response)
}
