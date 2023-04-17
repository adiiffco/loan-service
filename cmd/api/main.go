package api

import (
	"loanapp/middlewares/cors"
	"loanapp/middlewares/logger"
	"loanapp/middlewares/request"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func useGin() *gin.Engine {
	router := gin.New()
	router.Use(
		cors.Initialize(),
		gin.Recovery(),
		requestid.New(),
		request.AddRequestID(),
		logger.Initialize(),
	)
	return router
}

func Initialize() *gin.Engine {
	return useGin()
}
