package http

import (
	"auth-microservice/pkg/endpoint"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
)

// definition of encoder and decoder
func DecodeNewChallengeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.NewChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeVerifyChallengeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.VerifyChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func NewHTTPServer(ctx context.Context, endpoints endpoint.Endpoints) *gin.Engine {
	router := gin.Default()
	// use middlewares
	router.Use(commonMiddleware())
	auth := router.Group("/auth")
	{
		chap := auth.Group("/chap")
		{
			v1 := chap.Group("/v1")
			// needs to wrap http.Handler
			v1.POST("/new-challenge", gin.WrapH(httptransport.NewServer(
				endpoints.NewChallengeEndpoint,
				DecodeNewChallengeRequest,
				EncodeResponse,
			)))

			v1.POST("/verify-challenge", gin.WrapH(httptransport.NewServer(
				endpoints.VerifyChallengeEndpoint,
				DecodeVerifyChallengeRequest,
				EncodeResponse,
			)))
		}
	}
	return router
}

func commonMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Content-Type", "application/json")
	}
}
