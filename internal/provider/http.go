package provider

import (
	"github.com/JordanMarcelino/go-gin-starter/internal/config"
	"github.com/JordanMarcelino/go-gin-starter/internal/controller"
	"github.com/gin-gonic/gin"
)

func BootstrapHttp(cfg *config.Config, router *gin.Engine) {

	appController := controller.NewAppController()
	appController.Route(router)
}
