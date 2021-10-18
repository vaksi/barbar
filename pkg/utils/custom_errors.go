package utils

import (
	"fmt"
	"github.com/go-chi/render"
	"net/http"
)

type RequestError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Err        error  `json:"err,omitempty"`
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("%v", r.Err)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	if _, ok := err.(*RequestError); ok {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, RequestError{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, RequestError{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
	}
}
