package ginutils

import (
	"net/http"

	"github.com/JordanMarcelino/go-gin-starter/internal/constant"
	"github.com/JordanMarcelino/go-gin-starter/internal/dto"
	"github.com/gin-gonic/gin"
)

func ResponseOK[T any](ctx *gin.Context, data T) {
	ResponseJSON(ctx, http.StatusOK, constant.ResponseSuccessMessage, data, nil)
}

func ResponseOKPlain(ctx *gin.Context) {
	ResponseOK[any](ctx, nil)
}

func ResponseOKPagination[T any](ctx *gin.Context, data T, paging *dto.PageMetaData) {
	ResponseJSON(ctx, http.StatusOK, constant.ResponseSuccessMessage, data, paging)
}

func ResponseCreated[T any](ctx *gin.Context, data T) {
	ResponseJSON(ctx, http.StatusCreated, constant.ResponseSuccessMessage, data, nil)
}

func ResponseCreatedPlain(ctx *gin.Context) {
	ResponseCreated[any](ctx, nil)
}

func ResponseJSON[T any](ctx *gin.Context, statusCode int, message string, data T, paging *dto.PageMetaData) {
	ctx.JSON(statusCode, dto.WebResponse[T]{
		Message: message,
		Data:    data,
		Paging:  paging,
	})
}
