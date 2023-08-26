package notification

import (
	"album-manager/src/configs"
	"fmt"
	"net/smtp"
)

type Notifier interface {
	Send(msg string, to []string)
}

type EmailNotifier struct{}

func (EmailNotifier) Send(msg string, to []string) {
	from := configs.Env.Email.PrimaryEmail
	password := configs.Env.Email.Password
	host := configs.Env.Email.Host
	port := configs.Env.Email.Port
	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", host, port), auth, from, to, []byte(msg))
	if err != nil {
		fmt.Printf("Send email error: %v", err)
		return
	}

	fmt.Println("Send Email Successfully")
}

type SmsNotifier struct{}

func (SmsNotifier) Send(msg string, to []string) {
}

type Service struct {
	Notifier Notifier
}

func (s Service) SendNotification(msg string, to []string) {
	s.Notifier.Send(msg, to)
}

func CreateNotifier(t string) Notifier {
	if t == "sms" {
		return SmsNotifier{}
	}

	return EmailNotifier{}
}
