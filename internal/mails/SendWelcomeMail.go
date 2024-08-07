package mails

import (
	"log"

	"github.com/ubaniIsaac/go-project-manager/internal/helpers"
)

func SendWelcomeMail(
	recipient string,
	subject string,
	name string) error {
	values := struct {
		Name string
	}{
		Name: name,
	}

	templateFile := "../../internal/templates/welcome.html"

	err := helpers.DeliverMail(templateFile, values, recipient, subject)
	if err != nil {
		log.Fatalf("Failed to deliver mail: %v", err)
	}

	return nil
}
