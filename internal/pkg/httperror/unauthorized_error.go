package httperror

import (
	"errors"
	"net/http"

	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/constant"
)

func NewUnauthorizedError() *ResponseError {
	msg := constant.UnauthorizedErrorMessage

	err := errors.New(msg)

	return NewResponseError(err, http.StatusUnauthorized, msg)
}
