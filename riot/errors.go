package riot

import "fmt"

type RiotError struct {
	StatusCode int
	Message    string
}

func (e *RiotError) Error() string {
	return fmt.Sprintf("riot api error %d: %s", e.StatusCode, e.Message)
}

func (e *RiotError) IsNotFound() bool {
	return e.StatusCode == 404
}

func (e *RiotError) IsRateLimited() bool {
	return e.StatusCode == 429
}

func (e *RiotError) IsUnauthorized() bool {
	return e.StatusCode == 401
}

func (e *RiotError) IsForbidden() bool {
	return e.StatusCode == 403
}

func (e *RiotError) IsServerError() bool {
	return e.StatusCode >= 500
}

func NewRiotError(statusCode int, message string) *RiotError {
	return &RiotError{
		StatusCode: statusCode,
		Message:    message,
	}
}