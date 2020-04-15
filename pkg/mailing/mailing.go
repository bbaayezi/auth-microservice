package mailing

import (
	"net/smtp"
	"strings"
)

var auth = smtp.PlainAuth("Fox Mail", "robin_scrlt@foxmail.com", "scwuyrcsaqsmbjjg", "smtp.qq.com")

func SendMail(to []string, content string) error {
	// TODO: read configuration from envionment variable
	// to := []string{"476296987@qq.com"}
	nickname := "System"
	user := "robin_scrlt@foxmail.com"
	subject := "Verification Code"
	contentType := "Content-Type: text/plain; charset=UTF-8"
	address := "smtp.qq.com:25"
	// body := "This is the email body."
	msg := generateMailContent(to, user, nickname, subject, contentType, content)
	err := smtp.SendMail(address, auth, user, to, msg)
	return err
}

func generateMailContent(to []string, from string, nickname string, subject string, contentType string, body string) []byte {
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + from + ">\r\nSubject: " + subject + "\r\n" +
		contentType + "\r\n\r\n" + body)
	return msg
}
