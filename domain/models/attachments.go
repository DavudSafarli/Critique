package models

// Attachment is a models model
type Attachment struct {
	ID         uint   `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Path       string `json:"path,omitempty"`
	FeedbackID uint   `json:"feedback_id,omitempty"`
}
