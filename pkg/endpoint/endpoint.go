package endpoint

import (
	"auth-microservice/pkg/service"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	NewChallengeEndpoint    endpoint.Endpoint
	VerifyChallengeEndpoint endpoint.Endpoint
}

// define request and responses
type NewChallengeRequest struct {
	Key string `json:"key"`
}

type NewChallengeResponse struct {
	Challenge string `json:"challenge"`
	Err       string `json:"error,omitempty"`
}

type VerifyChallengeRequest struct {
	Key    string `json:"key"`
	Answer string `json:"answer"`
	Field  string `json:"field"`
}

type VerifyChallengeResponse struct {
	Correct bool   `json:"correct"`
	Err     string `json:"error,omitempty"`
}

// functions to create endpoints
func MakeNewChallengeEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(NewChallengeRequest)
		v, err := svc.NewChallenge(req.Key)
		if err != nil {
			return NewChallengeResponse{"", err.Error()}, nil
		}
		return NewChallengeResponse{v, ""}, nil
	}
}

func MakeVerifyChallengeEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(VerifyChallengeRequest)
		v, err := svc.VerifyChallenge(req.Key, req.Answer, req.Field)
		if err != nil {
			return VerifyChallengeResponse{false, err.Error()}, nil
		}
		return VerifyChallengeResponse{v, ""}, nil
	}
}
