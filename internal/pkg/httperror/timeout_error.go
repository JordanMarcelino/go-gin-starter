package httperror

import (
	"errors"
	"net/http"

	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/constant"
)

func NewTimeoutError() *ResponseError {
	msg := constant.RequestTimeoutErrorMessage

	err := errors.New(msg)

	return NewResponseError(err, http.StatusRequestTimeout, msg)
}
