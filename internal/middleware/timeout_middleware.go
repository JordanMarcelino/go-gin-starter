package middleware

import (
	"context"
	"time"

	"github.com/JordanMarcelino/go-gin-starter/internal/config"
	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/httperror"
	"github.com/gin-gonic/gin"
)

func RequestTimeout(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		timeoutCtx, cancel := context.WithTimeout(
			ctx.Request.Context(),
			time.Duration(cfg.HttpServer.RequestTimeoutPeriod)*time.Second,
		)
		defer cancel()

		done := make(chan struct{})
		ctx.Request = ctx.Request.WithContext(timeoutCtx)

		go next(ctx, done)

		select {
		case <-timeoutCtx.Done():
			ctx.Error(httperror.NewTimeoutError())
			ctx.Abort()
		case <-done:
		}
	}
}

func next(ctx *gin.Context, done chan struct{}) {
	defer func() {
		close(done)
		if err, ok := recover().(error); ok && err != nil {
			ctx.Error(err)
			ctx.Abort()
		}
	}()

	ctx.Next()
}
