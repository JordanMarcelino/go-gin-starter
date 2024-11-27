package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/JordanMarcelino/go-gin-starter/internal/config"
	"github.com/JordanMarcelino/go-gin-starter/internal/middleware"
	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/logger"
	"github.com/JordanMarcelino/go-gin-starter/internal/pkg/utils/validationutils"
	"github.com/JordanMarcelino/go-gin-starter/internal/provider"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type HttpServer struct {
	cfg    *config.Config
	server *http.Server
}

func NewHttpServer(cfg *config.Config) *HttpServer {
	gin.SetMode(cfg.App.Environment)

	router := gin.New()
	router.ContextWithFallback = true
	router.HandleMethodNotAllowed = true

	RegisterValidators()
	RegisterMiddleware(router, cfg)

	provider.BootstrapHttp(cfg, router)

	return &HttpServer{
		cfg: cfg,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.HttpServer.Host, cfg.HttpServer.Port),
			Handler: router,
		},
	}
}

func (s *HttpServer) Start() {
	logger.Log.Info("Running HTTP server on port:", s.cfg.HttpServer.Port)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Log.Fatal("Error while HTTP server listening:", err)
	}
	logger.Log.Info("HTTP server is not receiving new requests...")
}

func (s *HttpServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.cfg.HttpServer.GracePeriod)*time.Second)
	defer cancel()

	logger.Log.Info("Attempting to shut down the HTTP server...")
	if err := s.server.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Error shutting down HTTP server:", err)
	}
	logger.Log.Info("HTTP server shut down gracefully")
}

func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(validationutils.TagNameFormatter)
		v.RegisterCustomTypeFunc(validationutils.DecimalType)
	}
}

func RegisterMiddleware(router *gin.Engine, cfg *config.Config) {
	middlewares := []gin.HandlerFunc{
		gzip.Gzip(gzip.BestSpeed),
		middleware.Logger(),
		middleware.ErrorHandler(),
		middleware.RequestTimeout(cfg),
		cors.New(cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
			AllowAllOrigins:  true,
			AllowCredentials: true,
		}),
		gin.Recovery(),
	}

	router.Use(middlewares...)
}
