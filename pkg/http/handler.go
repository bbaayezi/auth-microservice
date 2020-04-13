package http

import (
	"auth-microservice/pkg/endpoint"
	"context"
	"encoding/json"
	"net/http"
)

// definition of encoder and decoder
func decodeNewChallengeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.NewChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeVerifyChallengeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request endpoint.VerifyChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
