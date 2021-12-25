package main

import (
	"github.com/VishalTanwani/gochat-apiserver/internal/models"
	"github.com/xhit/go-simple-mail/v2"
	"time"
	"fmt"
	// "flag"
)

var emailID string
var emailPass string

func setPass(email,pass *string){
	emailID = *email
	emailPass = *pass
}

func listenForMail(){
	for {
		msg := <- app.MailChan
		sendMail(msg)
	}
}

func sendMail(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587
	server.Username = emailID
	server.Password = emailPass
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client,err := server.Connect()
	if err!=nil {
		fmt.Println("error at connecting mail server",err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)
	
	err = email.Send(client)
	if err != nil {
		fmt.Println("error at sending email",err)
	} else {
		fmt.Println("MailSend")
	}
}