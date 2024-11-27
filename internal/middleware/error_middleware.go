package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/JordanMarcelino/go-gin-starter/internal/constant"
	"github.com/JordanMarcelino/go-gin-starter/internal/dto"
	pkgconstant "github.com/JordanMarcelino/go-gin-starter/internal/pkg/constant"
	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/httperror"
	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/utils/validationutils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		errLen := len(ctx.Errors)
		if errLen > 0 {
			err := ctx.Errors.Last()

			switch e := err.Err.(type) {
			case validator.ValidationErrors:
				handleValidationError(ctx, e)
			case *json.SyntaxError:
				handleJsonSyntaxError(ctx)
			case *json.UnmarshalTypeError:
				handleJsonUnmarshalTypeError(ctx, e)
			case *time.ParseError:
				handleParseTimeError(ctx, e)
			case *httperror.ResponseError:
				ctx.AbortWithStatusJSON(e.GetCode(), dto.WebResponse[any]{
					Message: e.DisplayMessage(),
				})
			default:
				if errors.Is(e, io.EOF) {
					ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
						Message: pkgconstant.EOFErrorMessage,
					})
					return
				}

				ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.WebResponse[any]{
					Message: pkgconstant.InternalServerErrorMessage,
				})
			}
		}
	}
}

func handleJsonSyntaxError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: pkgconstant.JsonSyntaxErrorMessage,
	})
}

func handleJsonUnmarshalTypeError(ctx *gin.Context, err *json.UnmarshalTypeError) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: fmt.Sprintf(pkgconstant.JsonUnmarshallTypeErrorMessage, err.Field),
	})
}

func handleParseTimeError(ctx *gin.Context, err *time.ParseError) {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: fmt.Sprintf("please send time in format of %s, got: %s", constant.ConvertGoTimeLayoutToReadable(err.Layout), err.Value),
	})
}

func handleValidationError(ctx *gin.Context, err validator.ValidationErrors) {
	ve := []dto.FieldError{}

	for _, fe := range err {
		ve = append(ve, dto.FieldError{
			Field:   fe.Field(),
			Message: validationutils.TagToMsg(fe),
		})
	}

	ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.WebResponse[any]{
		Message: pkgconstant.ValidationErrorMessage,
		Errors:  ve,
	})
}
