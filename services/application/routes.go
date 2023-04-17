package application

import (
	"loanapp/middlewares/auth"

	"github.com/gin-gonic/gin"
)

var appController *Application

func internalRouteHandler(router *gin.RouterGroup) {
	router.Use(auth.Authorize())
	router.GET("/metadata", appController.MetaData)
	router.POST("/submit", appController.SubmitDetails)
	router.POST("/verify", appController.VerifyDetails)
	router.GET("/balance-sheet", appController.BalanceSheet)
	router.POST("/decision", appController.Decision)
}

func externalRouteHandler(router *gin.RouterGroup) {
	router.GET("/authorize", appController.Authorize)
}

func RouteHandler(router *gin.RouterGroup) {
	appController = InitController()
	internalRouteHandler(router.Group("/v1/internal"))
	externalRouteHandler(router.Group("/v1/external"))
}
