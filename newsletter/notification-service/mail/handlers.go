package mail

import "log"

type MailPayload struct {
	ID        int
	FirstName string `bson:"first_name" json:"first_name"`
	Name      string `bson:"name" json:"name"`
	To        string `bson:"to" json:"to"`
	Subject   string `bson:"subject" json:"subject"`
	Message   string `bson:"message" json:"message"`
}

// This handler is responsible for calling the required function for sending mail.
func SendMail(mail MailPayload) error {
	msg := Message{
		FirstName: mail.FirstName,
		From:      "",
		To:        mail.To,
		Subject:   mail.Subject,
		Data:      mail.Message,
	}

	err := SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
