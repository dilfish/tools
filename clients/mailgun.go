// Copyright 2018 Sean.ZH

package clients

import (
	"context"
	"time"

	mailgun "github.com/mailgun/mailgun-go/v4"
)

// ApiKey of mailgun
var ApiKey = ""

// InitMail set api and pub key
func InitMail(api string) {
	ApiKey = api
}

// SendSimpleMessage send email via mailgun
func SendSimpleMessage(title, content, to string) (string, error) {
	domain := "dilfish.icu"
	apiKey := ApiKey
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
		"Mc Noticer<mcnotice@dilfish.icu>",
		title,
		content,
		to,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, m)
	return id, err
}
