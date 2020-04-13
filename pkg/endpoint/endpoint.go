package endpoint

import (
	"auth-microservice/pkg/service"
	"context"

	"github.com/go-kit/kit/endpoint"
)

// define request and responses
type newChallengeRequest struct {
	Key string `json:"key"`
}

type newChallengeResponse struct {
	Challenge string `json:"challenge"`
	Err       string `json:"error,omitempty"`
}

type verifyChallengeRequest struct {
	Key    string `json:"key"`
	Answer string `json:"answer"`
	Field  string `json:"field"`
}

type verifyChallengeResponse struct {
	Correct bool   `json:"correct"`
	Err     string `json:"error,omitempty"`
}

// functions to create endpoints
func makeNewChallengeEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(newChallengeRequest)
		v, err := svc.NewChallenge(req.Key)
		if err != nil {
			return newChallengeResponse{"", err.Error()}, nil
		}
		return newChallengeResponse{v, ""}, nil
	}
}

func makeVerifyChallengeEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(verifyChallengeRequest)
		v, err := svc.VerifyChallenge(req.Key, req.Answer, req.Field)
		if err != nil {
			return verifyChallengeResponse{false, err.Error()}, nil
		}
		return verifyChallengeResponse{v, ""}, nil
	}
}
