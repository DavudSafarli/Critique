package domain

import (
	"errors"
)

var (
	// ErrNotFound happens when requested resource does not exist. Ex: in storage
	ErrNotFound = errors.New("Requested resource not found")
	// ErrUnprocessable happens when user input is unable to be processed. Ex: can't be decoded, or unvalid input
	ErrUnprocessable = errors.New("Input is unprocessable")
	// ErrInternalServerError happens in all other cases :/
	ErrInternalServerError = errors.New("ErrInternalServerError")
)
