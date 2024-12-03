package utils

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"os"
	"strconv"
)

func SendEmail(to string, subject string, data interface{}, templateFile string) error {
	from := "no-reply@kickof.com"

	result, _ := ParseTemplate(templateFile, data)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", result)

	port, _ := strconv.ParseInt(os.Getenv("MAIL_PORT"), 10, 64)
	d := gomail.NewDialer(os.Getenv("MAIL_HOST"), int(port), os.Getenv("MAIL_EMAIL"), os.Getenv("MAIL_APP_PASSWORD"))
	err := d.DialAndSend(m)
	if err != nil {
		panic(err)
	}
	return err
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		log.Println(err)
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		log.Println(err)
		return "", err
	}
	return buf.String(), nil
}

func SendEmailVerification(template string, to string, data interface{}) (bool, error) {
	var err error
	t := "templates/" + template + ".html"
	subject := "Please Confirm Your Email"
	err = SendEmail(to, subject, data, t)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}
