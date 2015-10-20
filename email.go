package ice

import (
	"github.com/mailgun/mailgun-go"
	"log"
)

func SendMail(to, subject, body string) {
	mg := mailgun.NewMailgun(Config.MailgunDomain, Config.MailgunKey, "")
	msg := mg.NewMessage(Config.EmailFrom, subject, body, to)
	mes, id, err := mg.Send(msg)
	log.Println(mes, id, err)
}
