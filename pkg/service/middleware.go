package service

import (
	"time"

	"github.com/go-kit/kit/log"
)

// service level middleware
type loggingMiddleware struct {
	logger log.Logger
	next   AuthService
}

func (mw loggingMiddleware) NewChallenge(key string) (challenge string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "newChallenge",
			"key", key,
			"output", challenge,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	challenge, err = mw.next.NewChallenge(key)
	return
}
func (mw loggingMiddleware) VerifyChallenge(key string, answer string, field string) (correct bool, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "verifyChallenge",
			"key", key,
			"client_answer", answer,
			"field_string", field,
			"correct", correct,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	correct, err = mw.next.VerifyChallenge(key, answer, field)
	return
}
func (mw loggingMiddleware) Salting(str string) (salt string, result string) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "salting",
			"input", str,
			"salt", salt,
			"output", result,
			"took", time.Since(begin),
		)
	}(time.Now())
	salt, result = mw.next.Salting(str)
	return
}
func (mw loggingMiddleware) SendToken(ct ContactService) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "sendToken",
			"contact_id", ct.GetContactID(),
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.next.SendToken(ct)
	return
}
func (mw loggingMiddleware) VerifyToken(token string, ct ContactService) (correct bool, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "verifyToken",
			"contact_id", ct.GetContactID(),
			"client_token", token,
			"correct", correct,
			"error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	correct, err = mw.next.VerifyToken(token, ct)
	return
}
