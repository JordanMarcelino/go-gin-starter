package middleware

import (
	"net/http"
	"time"

	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/httperror"
	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		ctx.Next()

		params := map[string]any{
			"status_code": ctx.Writer.Status(),
			"client_ip":   ctx.ClientIP(),
			"method":      ctx.Request.Method,
			"latency":     time.Since(start).String(),
			"path":        path,
		}

		if len(ctx.Errors) == 0 {
			logger.Log.WithFields(params).Info("incoming request")
			return
		}
		logErrors(ctx, params)
	}
}

func logErrors(ctx *gin.Context, params map[string]any) {
	errors := []error{}
	for _, err := range ctx.Errors {
		switch e := err.Err.(type) {
		case validator.ValidationErrors:
			params["status_code"] = http.StatusBadRequest
			errors = append(errors, err)
		case *httperror.ResponseError:
			params["status_code"] = e.GetCode()
			errors = append(errors, e.OriginalError())
		default:
			params["status_code"] = http.StatusInternalServerError
			errors = append(errors, err)
		}
	}

	params["errors"] = errors
	logger.Log.WithFields(params).Error("got error")
}
