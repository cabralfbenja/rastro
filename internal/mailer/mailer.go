package mailer

import (
	"bytes"
	"embed"
	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
)

//go:embed templates
var templateFS embed.FS

type Mailer struct {
	dialer *gomail.Dialer
	sender string // no-reply@rastro.com
	Logger echo.Logger
}

type EmailData struct {
	AppName string
	Subject string
	Meta    interface{}
}

func NewMailer(logger echo.Logger) *Mailer {
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		logger.Fatal(err)
	}

	mailHost := os.Getenv("MAIL_HOST")
	mailUser := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailSender := os.Getenv("MAIL_SENDER")

	dialer := gomail.NewDialer(mailHost, mailPort, mailUser, mailPassword)

	return &Mailer{
		dialer: dialer,
		sender: mailSender,
		Logger: logger,
	}
}

func (mailer *Mailer) Send(recipient string, templateFile string, data EmailData) error {
	absPath := filepath.Join("templates", templateFile)
	tmpl, err := template.ParseFS(templateFS, absPath)

	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	data.AppName = os.Getenv("APP_NAME")

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}

	gomailMessage := gomail.NewMessage()
	gomailMessage.SetHeader("From", mailer.sender)
	gomailMessage.SetHeader("To", recipient)
	gomailMessage.SetHeader("Subject", subject.String())

	gomailMessage.SetBody("text/html", htmlBody.String())

	err = mailer.dialer.DialAndSend(gomailMessage)
	if err != nil {
		mailer.Logger.Error(err)
		return err
	}
	return nil
}
