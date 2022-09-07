package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"sync"
)

var (
	From_mail     = os.Getenv("FROM_MAIL")
	Mail_password = os.Getenv("MAIL_PASSWORD")
	SMTP_Host     = os.Getenv("HOST")
	Mail_subject  string
	Mail_body     string

	from      *mail.Address
	auth      smtp.Auth
	tlsconfig *tls.Config
	mailwg    sync.WaitGroup
)

type Container struct {
	m       sync.Mutex
	Headers map[string]string
}

func NewContainer() *Container {
	return &Container{
		Headers: make(map[string]string),
	}
}

func init() {
	from = &mail.Address{Name: "Test-mail", Address: From_mail}
	auth = smtp.PlainAuth("", From_mail, Mail_password, SMTP_Host)
	tlsconfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         SMTP_Host,
	}
}

func SendSSLMail(subject, msg string, recipient string) {
	to := mail.Address{Name: "", Address: recipient}

	Mail_subject = subject
	Mail_body = msg

	// initialize new container object
	container := NewContainer()
	// call mutex.lock to avoid multiple writes to
	// one header instance from running goroutines
	container.m.Lock()
	container.Headers["From"] = from.String()
	container.Headers["To"] = to.String()
	container.Headers["Subject"] = Mail_subject
	// unlock mutex after function returns
	defer container.m.Unlock()

	// Setup message
	message := ""
	for k, v := range container.Headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + Mail_body

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", SMTP_Host, 465), tlsconfig)
	if err != nil {
		log.Printf("Error sending mail %v", err)
		return
	}

	c, err := smtp.NewClient(conn, SMTP_Host)
	if err != nil {
		log.Printf("Error sending mail %v", err)
		return
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Printf("Error sending mail %v", err)
		return
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Printf("Error sending mail %v", err)
		return
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Printf("Error sending mail %v", err)
		return
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Printf("Error sending mail %v", err)
		return
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Printf("Error sending mail %v", err)
		return
	}

	err = w.Close()
	if err != nil {
		log.Printf("Error sending mail %v", err)
		return
	}

	if err = c.Quit(); err != nil {
		return
	}
}

// Concurrently sending mails to multiple recipients
func SendMails(subject, msg string, recipients []string) {
	mailwg.Add(len(recipients))
	for _, v := range recipients {
		go func(recipient string) {
			defer mailwg.Done()
			SendSSLMail(subject, msg, recipient)
		}(v)
	}
	mailwg.Wait()
}
