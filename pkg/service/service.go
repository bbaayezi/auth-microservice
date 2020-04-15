package service

import (
	"auth-microservice/pkg/db"
	"auth-microservice/pkg/mailing"
	"auth-microservice/pkg/utils"
	"crypto/md5"
	"fmt"
	"io"
	"time"

	"github.com/go-kit/kit/log"
)

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
	Salting(str string) (salt string, result string)
}

type TokenBaseService interface {
	SendToken(ct ContactService) error
	VerifyToken(token string, id string) (correct bool, err error)
}

type ContactService interface {
	GetContactID() string
	SendMessage(msg string) error
}

// implementing services
type authService struct{}

func (authService) NewChallenge(key string) (challenge string, err error) {
	// get redis connection client
	redisClient, err := db.GetRedisClient()
	if err != nil {
		return "", err
	}
	// generate a new 8 bit long challenge
	c := utils.GenerateNonce(8)
	// set key if not exist
	setReq := redisClient.SetNX(key, c, 3*time.Minute)
	setSuccess, err := setReq.Result()
	if !setSuccess {
		// returns the current unsolved challenge
		challenge, err = redisClient.Get(key).Result()
	} else {
		challenge, err = c, nil
	}
	return
}

func (authService) VerifyChallenge(key string, answer string, field string) (correct bool, err error) {
	// get key from redis
	redisClient, err := db.GetRedisClient()
	if err != nil {
		return false, err
	}
	c, err := redisClient.Get(key).Result()
	if err != nil {
		return false, err
	}
	// calculate the challenge answer
	hash := md5.New()
	io.WriteString(hash, field+c)
	ans := fmt.Sprintf("%x", hash.Sum(nil))
	// compare the answer
	correct = ans == answer
	if correct {
		// delete challenge
		_, err = redisClient.Del(key).Result()
	}
	return
}

func (authService) Salting(str string) (salt string, result string) {
	// generate a salt
	salt = utils.GenerateNonce(6)
	// hashing
	hash := md5.New()
	io.WriteString(hash, str+salt)
	result = fmt.Sprintf("%x", hash.Sum(nil))
	return
}

func (authService) SendToken(ct ContactService) error {
	// get redis client to store token
	redisClient, err := db.GetRedisClient()
	if err != nil {
		return err
	}
	id := ct.GetContactID()
	token := utils.GenerateNonce(8)
	key := id + ".token"
	_, err = redisClient.SetNX(key, token, 5*time.Minute).Result()
	if err != nil {
		return err
	}
	// send token
	return ct.SendMessage(fmt.Sprintf("Your verification code is: %s, it will expire in 5 minutes.", token))
}

func (authService) VerifyToken(token string, id string) (correct bool, err error) {
	// get redis client to store token
	redisClient, err := db.GetRedisClient()
	if err != nil {
		return false, err
	}
	key := id + ".token"
	t, err := redisClient.Get(key).Result()
	if err != nil {
		return false, err
	}
	// compare token
	correct = t == token
	if correct {
		// delete token
		_, err = redisClient.Del(key).Result()
	}
	return
}

func NewAuthService(logger log.Logger) AuthService {
	var svc AuthService
	svc = authService{}
	// apply service level middleware
	svc = loggingMiddleware{
		logger,
		svc,
	}
	return svc
}

type EmailContact struct {
	Email string
	ID    string
}

func (ect EmailContact) GetContactID() string {
	return ect.ID
}

func (ect EmailContact) SendMessage(msg string) error {
	go mailing.SendMail([]string{ect.Email}, msg)
	return nil
}
