package web

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

var (
	ErrUnsuccessfulAction = errors.New("unsuccessful action")
	ErrLoginRequired      = errors.New("login required")
)

// Size amd64: 8+(8*2)+(8*2) = 45 bytes
type ErrInvalidResponse struct {
	resp *http.Response
	msg  string
	err  error
}

// Response return http.Response with [`Body`] closed.
func (e ErrInvalidResponse) Response() *http.Response {
	return e.resp
}

func (e ErrInvalidResponse) Unwrap() error {
	return e.err
}

func (e ErrInvalidResponse) Error() string {
	if e.msg == "" {
		return "invalid response"
	}
	return fmt.Sprintf("invalid response: %s", e.msg)
}
