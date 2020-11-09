package emails

import (
	"backend/types"
	log "github.com/sirupsen/logrus"
	"net/smtp"
	"os"
)

type EmailManager struct {
}

var Manager types.EmailService = &EmailManager{}

func (m *EmailManager) SendVerificationKey(to, key string) {
	//link := "https://aim-love.ga/verify/" + key
	link := "http://localhost:" + os.Getenv("BACKEND_PORT") + "/api/v1/verify/" + key
	body := `<h3>Hello from Matcha!</h3><p>To verify this email address, follow this <a href="` + link + `">link</a></p>`
	m.sendEmailFromService(to, "Verify email", body)
}

func (m *EmailManager) SendGoodbyeMessage(to string) {
	body := `<h3>Goodbye from Matcha!</h3><p>Your account has been successfully deleted! <br>Good luck!</p>`

	m.sendEmailFromService(to, "Good bye!", body)
}

func (m *EmailManager) sendEmailFromService(to, subject, body string) {
	from := os.Getenv("SERVICE_MAIL_ADDR")
	pass := os.Getenv("SERVICE_MAIL_PASSWD")
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" + mime + body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Errorf("smtp error: %s", err)
	}
}