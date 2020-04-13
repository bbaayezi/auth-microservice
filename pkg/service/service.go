package service

// AuthService is the combination of serveral specific authentication methods
type AuthService interface {
	CHAPService
	PAPService
	TokenBaseService
}

type CHAPService interface {
	// NewChallenge creates a challenge for the specified key
	NewChallenge(key string) (challenge string, err error)
	// VerifyChallenge verifies the challenge answer for the specified key provided with field string
	VerifyChallenge(key string, answer string, field string) (correct bool, err error)
}

type PAPService interface {
	Salting(str string) (salt string, err error)
}

type TokenBaseService interface {
	SendToken(ct ContactService) error
	VerifyToken(token string, ct ContactService) (correct bool, err error)
}

type ContactService interface {
	GetContactID() string
	SendMessage(msg string) error
}
