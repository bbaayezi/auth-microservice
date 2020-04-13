package service

import (
	"auth-microservice/pkg/db"
	"auth-microservice/pkg/utils"
	"crypto/md5"
	"fmt"
	"io"
	"time"
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
	VerifyToken(token string, ct ContactService) (correct bool, err error)
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
