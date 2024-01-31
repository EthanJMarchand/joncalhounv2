package main

import (
	"fmt"

	"github.com/ethanjmachand/lenslocked/models"
)

const (
	host     = "sandbox.smtp.mailtrap.io"
	port     = 587
	username = "c2774a8554d065"
	password = "ce6b1672d58043"
)

func main() {
	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})
	err := es.ForgotPassword("ethan@k2r.ca", "https://lenslocked.com/reset-pw?token=abc123")
	if err != nil {
		panic(err)
	}

	fmt.Println("Message sent!")
}
