package main

import (
	"auth-microservice/pkg/endpoint"
	"auth-microservice/pkg/http"
	"auth-microservice/pkg/service"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	authSvc := service.NewAuthService()

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
}
