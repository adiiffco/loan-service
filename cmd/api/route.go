package api

import (
	"context"
	"fmt"
	"net/http"

	"loanapp/services/application"
	"loanapp/services/health"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var (
	apiEngine *gin.Engine
	srv       *http.Server
)

func Run(httpShutdown <-chan bool) {
	port := fmt.Sprintf(":%s", viper.GetString("port"))
	srv = &http.Server{
		Addr:    port,
		Handler: apiEngine,
	}
	fmt.Println("service running on port", port)
	srv.ListenAndServe()
	<-httpShutdown
	fmt.Println("service shutting down")
}

func ShutDown(ctx context.Context) {
	if err := srv.Shutdown(ctx); err != nil {
		panic(fmt.Errorf("can't shutdown server, err: %s", err))
	}
}

func InitRoutes() {
	fmt.Println("Initializing routes")
	apiEngine = Initialize()
	loanRoutes := apiEngine.Group("/loan")
	health.RouteHandler(apiEngine.Group("/health"))
	application.RouteHandler(loanRoutes.Group("/application"))
}
