package endpoint

import (
	"auth-microservice/pkg/service"
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	NewChallengeEndpoint    endpoint.Endpoint
	VerifyChallengeEndpoint endpoint.Endpoint
	SaltingEndpoint         endpoint.Endpoint
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

// PAP service request/response
type SaltingRequest struct {
	Str string `json:"string"`
}

type SaltingResponse struct {
	Salt   string `json:"salt"`
	Result string `json:"result"`
}

// Token based service
type SendTokenRequest struct {
	ContactType string `json:"contactType"`
	Message     string `json:"message"`
}

type SendTokenResponse struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
	Err     string `json:"error,omitempty"`
}

type VerifyTokenRequest struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type VerifyTokenResponse struct {
	Correct bool   `json:"string"`
	Err     string `json:"error,omitempty"`
}

// functions to create endpoints
func makeNewChallengeEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(NewChallengeRequest)
		v, err := svc.NewChallenge(req.Key)
		if err != nil {
			return NewChallengeResponse{"", err.Error()}, nil
		}
		return NewChallengeResponse{v, ""}, nil
	}
}

func makeVerifyChallengeEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(VerifyChallengeRequest)
		v, err := svc.VerifyChallenge(req.Key, req.Answer, req.Field)
		if err != nil {
			return VerifyChallengeResponse{false, err.Error()}, nil
		}
		return VerifyChallengeResponse{v, ""}, nil
	}
}

// PAP service endpoint
func makeSaltingEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SaltingRequest)
		salt, result := svc.Salting(req.Str)
		return SaltingResponse{salt, result}, nil
	}
}

func New(svc service.AuthService) Endpoints {
	eps := Endpoints{
		NewChallengeEndpoint:    makeNewChallengeEndpoint(svc),
		VerifyChallengeEndpoint: makeVerifyChallengeEndpoint(svc),
		SaltingEndpoint:         makeSaltingEndpoint(svc),
	}
	// apply endpoint level middleware
	return eps
}
