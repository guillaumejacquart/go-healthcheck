package pkg

import (
	"fmt"
	"net/smtp"

	"github.com/guillaumejacquart/go-healthcheck/pkg/domain"
	"github.com/spf13/viper"
)

func sendNotification(a domain.App) error {
	if !a.Notify {
		return nil
	}

	fmt.Println("Sending mail for app ", a.Name)

	smtpHost := viper.GetString("smtp.host")
	smtpPort := viper.GetInt("smtp.port")
	smtpUsername := viper.GetString("smtp.username")
	smtpPassword := viper.GetString("smtp.password")
	smtpFrom := viper.GetString("smtp.from")
	smtpTo := viper.GetString("smtp.to")

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		smtpUsername,
		smtpPassword,
		smtpHost,
	)

	var content = fmt.Sprintf("Hello\nApp %v is %v.\ngo-healthcheck", a.Name, a.Status)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpHost+":"+fmt.Sprint(smtpPort),
		auth,
		smtpFrom,
		[]string{smtpTo},
		[]byte(content),
	)
    fmt.Println(err)
	return err
}