package tools

import (
	"fmt"
	"gopkg.in/mailgun/mailgun-go.v1"
	"errors"
)

var ApiKey = ""
var PubKey = ""

func InitMail(api, pub string) {
	ApiKey = api
	PubKey = pub
}

func SendMail(to, title, content string) error {
	if ApiKey == "" || PubKey == "" {
		return errors.New("You need to call InitMail first")
	}
	from := "mc@mg.libsm.com"
	domain := "mg.libsm.com"
	mg := mailgun.NewMailgun(domain, ApiKey, PubKey)
	m := mailgun.NewMessage(
		from, title,
		content, to)
	resp, id, err := mg.Send(m)
	fmt.Println("id", id, "resp", resp)
	return err
}