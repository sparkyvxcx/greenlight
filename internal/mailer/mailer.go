package mailer

import (
	"bytes"
	"embed"
	"text/template"
	"time"

	"github.com/go-mail/mail/v2"
)

// Declare a new variable with the type embed.FS (embedded file system) to hold email templates.
// Indicates to Go that we want to store the contents of the ./templates dir in the templateFS
// embedded file system variable.

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	// Initialize a new mail.Dialer instance with the given SMTP server settings.
	dialer := mail.NewDialer(host, port, username, password)
	// Set a 5-second timeout whenever we send an email out.
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer,
		sender,
	}
}

func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
	// Use the ParseFS() method to parse the required template file from the embedded file system.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Execute the named template "subject", passing in the dynamic data and storing the result in
	// a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Follow the same pattern to execute the "plainBody" template and store the result in the plainBody
	// variable.
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// And likewise with the "htmlBody" template.
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Use the mail.NewMessage() function to Initialize a new mail.Message instance.
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetHeader("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// Call the DialAndSend() method on the dialer, passing in the message to send.
	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
