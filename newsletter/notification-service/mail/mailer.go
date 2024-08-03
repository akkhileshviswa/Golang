package mail

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
)

type Message struct {
	FirstName string
	From      string
	To        string
	Subject   string
	Data      any
	DataMap   map[string]any
}

// This function is used to send the mail.
func SendSMTPMessage(msg Message) error {
	data := map[string]any{
		"name":    msg.FirstName,
		"message": msg.Data,
	}
	msg.DataMap = data

	plainMessage, err := buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	from := ""
	password := ""
	toEmailAddress := ""
	to := []string{toEmailAddress}
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	subject := "Subject: " + msg.Subject + "\n"
	body := plainMessage
	message := []byte(subject + body)
	auth := smtp.PlainAuth("", from, password, host)

	err = smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		log.Println(1)
		log.Println(err)
		return err
	}

	return nil
}

// This function is used to build plain text message using the template file,
// using the message recieved.
func buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./../../mail/templates/mail.plain.gohtml"
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
