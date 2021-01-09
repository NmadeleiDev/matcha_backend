package emails

import (
	"fmt"
	"net/smtp"
	"os"

	"backend/dao"

	log "github.com/sirupsen/logrus"
)

type EmailManager struct {
}

var Manager dao.EmailService = &EmailManager{}

var host = os.Getenv("PROJECT_HOST")

func (m *EmailManager) SendPasswordResetEmail(to, key string) {
	var link string
	if len(host) == 0 {
		link = "http://localhost:" + os.Getenv("PROJECT_PORT") + "/reset?k=" + key
	} else {
		link = fmt.Sprintf("https://%v/reset?k=%v", host, key)
	}

	template := `
<div>
    <h2>Password reset</h2>
    <p>If you didn't request password reset for your Matcha account, just ignore this email.</p>
    <p><a href="%v">Press here to reset password</a></p>
</div>`
	body := fmt.Sprintf(template, link)
	m.sendEmailFromService(to, "Matcha password reset", body)
}

func (m *EmailManager) SendAccountVerificationKey(to, key string) {
	var link string
	if len(host) == 0 {
		link = fmt.Sprintf("http://localhost:%v/api/main/verify/%v", os.Getenv("PROJECT_PORT"), key)
	} else {
		link = fmt.Sprintf("https://%v/api/main/verify/%v", host, key)
	}
	body := `<h3>Hello from Matcha!</h3><p>To verify this email address, follow this <a href="` + link + `">link</a></p>`
	m.sendEmailFromService(to, "Verify email", body)
}

func (m *EmailManager) SendEmailVerificationKey(to, key string)  {
	var link string
	if len(host) == 0 {
		link = fmt.Sprintf("http://localhost:%v/email/verify?key=%v", os.Getenv("PROJECT_PORT"), key)
	} else {
		link = fmt.Sprintf("https://%v/email/verify?key=%v", host, key)
	}
	body := `<h3>Change email</h3><p>To verify this new email address, follow this <a href="` + link + `">link</a></p>`
	m.sendEmailFromService(to, "Change email", body)
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
		"ToChat: " + to + "\n" +
		"Subject: " + subject + "\n" + mime + body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Errorf("smtp error: %s; addr=%v, passwd=%v", err, from, pass)
	}
	log.Infof("Email send: %v", msg)
}
