package endpoint

import (
	"auth-microservice/pkg/service"
	"context"

	"github.com/go-kit/kit/endpoint"
	uuid "github.com/satori/go.uuid"
)

type Endpoints struct {
	NewChallengeEndpoint    endpoint.Endpoint
	VerifyChallengeEndpoint endpoint.Endpoint
	SaltingEndpoint         endpoint.Endpoint
	SendTokenEndpoint       endpoint.Endpoint
	VerifyTokenEndpoint     endpoint.Endpoint
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
	ContactMethod string `json:"contactMethod"`
	Contact       string `json:"contact"`
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
	Correct bool   `json:"correct"`
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

// Token based service endpoints
func makeSendTokenEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SendTokenRequest)
		var id string
		// check contect method
		switch req.ContactMethod {
		case "email":
			// create a new email contact
			// generate a uuid
			id = uuid.Must(uuid.NewV4(), err).String()
			if err != nil {
				return SendTokenResponse{"", false, err.Error()}, nil
			}
			contact := service.EmailContact{
				// generate a uuid
				ID:    id,
				Email: req.Contact,
			}
			err := svc.SendToken(contact)
			if err != nil {
				return SendTokenResponse{"", false, err.Error()}, nil
			}
			break
		// other contact method goes here...
		default:
			return SendTokenResponse{"", false, "Unsupported contact method"}, nil
		}
		return SendTokenResponse{id, true, ""}, nil
	}
}

func makeVerifyTokenEndpoint(svc service.AuthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(VerifyTokenRequest)
		correct, err := svc.VerifyToken(req.Token, req.ID)
		if err != nil {
			return VerifyTokenResponse{false, err.Error()}, nil
		}
		return VerifyTokenResponse{correct, ""}, nil
	}
}

func New(svc service.AuthService) Endpoints {
	eps := Endpoints{
		NewChallengeEndpoint:    makeNewChallengeEndpoint(svc),
		VerifyChallengeEndpoint: makeVerifyChallengeEndpoint(svc),
		SaltingEndpoint:         makeSaltingEndpoint(svc),
		SendTokenEndpoint:       makeSendTokenEndpoint(svc),
		VerifyTokenEndpoint:     makeVerifyTokenEndpoint(svc),
	}
	// apply endpoint level middleware
	return eps
}
