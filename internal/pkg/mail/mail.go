package mail

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendMail(to string, subject string, body string) error
}

type mailer struct {
	from     string
	smtpHost string
	smtpPort int
	username string
	password string
}

func NewMailer(from, smtpHost string, smtpPort int, username, password string) Mailer {
	return &mailer{
		from:     from,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
	}
}

func (m *mailer) SendMail(to string, subject string, body string) error {
	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", m.from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)

	// Set email body
	message.SetBody("text/html", body)

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(m.smtpHost, m.smtpPort, m.username, m.password)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
	}
	return nil
}
