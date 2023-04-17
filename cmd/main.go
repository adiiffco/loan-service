package main

import (
	"context"
	"fmt"
	"loanapp/adapters/cache"
	"loanapp/adapters/logger"
	"loanapp/adapters/mysql"
	"loanapp/cmd/api"
	"loanapp/cmd/go-routines"
	c "loanapp/config"
	"loanapp/middlewares/signal"
	"time"
)

var httpShutdown = make(chan bool)

func init() {
	fmt.Println("Initializing service")
	ctx, cancel := context.WithCancel(context.Background())
	c.InitializeEnv()
	routines.InitializeMetrics()
	signal.SetupSignals()
	cache.Initialize(ctx)
	logger.Initialize()
	mysql.Initialize()
	api.InitRoutes()
	signal.CleanupOnSignal(func() {
		fmt.Println("Cleaning up....")
		cancel()
		api.ShutDown(context.Background())
		time.Sleep(time.Duration(5) * time.Second)
		routines.ShutDown()
		httpShutdown <- true
	})
}

func main() {
	api.Run(httpShutdown)
}
