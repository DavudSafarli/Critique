package api

import (
	"net/http"

	"github.com/DavudSafarli/Critique/domain"
)

// ErrorResponse represents the body of "not 200 responses"
type ErrorResponse struct {
	Status int
	//Type    string
	//Message string
}

func toHttp(err error) ErrorResponse {
	httpStatusCode := convertDomainErrorToHttpStatusCode(err)
	return ErrorResponse{
		Status: httpStatusCode,
	}
}

func convertDomainErrorToHttpStatusCode(domainError error) int {
	switch domainError {
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrUnprocessable:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
