package emails

import (
	"log"
	"net/smtp"
	"os"
)

func Send(to, key string) {
	from := os.Getenv("SERVICE_MAIL_ADDR")
	pass := os.Getenv("SERVICE_MAIL_PASSWD")

	//link := "https://aim-love.ga/verify/" + key
	link := "http://localhost:8081/api/v1/verify/" + key
	body := `<h3>Hello from Matcha!</h3><p>To verify this email address, follow this <a href="` + link + `">link</a></p>`
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"


	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Verify Email\n" + mime + body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent")
}
