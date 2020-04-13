package main

import (
	"auth-microservice/pkg/endpoint"
	"auth-microservice/pkg/service"
	"context"
	"fmt"
	"os/signal"
	"syscall"

	// "log"
	"auth-microservice/pkg/http"
	"os"

	"github.com/go-kit/kit/log"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stderr)
	ctx := context.Background()
	authSvc := service.NewAuthService(logger)

	// error handling
	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	// endpoints
	eps := endpoint.Endpoints{
		NewChallengeEndpoint:    endpoint.MakeNewChallengeEndpoint(authSvc),
		VerifyChallengeEndpoint: endpoint.MakeVerifyChallengeEndpoint(authSvc),
	}

	go func() {
		errChan <- http.NewHTTPServer(ctx, eps).Run(":8080")
	}()
	logger.Log("err: ", <-errChan)
}
