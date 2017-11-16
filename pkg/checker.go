package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	health "github.com/docker/go-healthcheck"
	"github.com/guillaumejacquart/go-healthcheck/pkg/domain"
	"github.com/spf13/viper"
)

func runChecksApp() {
	apps, err := getAllApps()

	if err != nil {
		panic(err)
	}

	for _, a := range apps {
		registerCheck(a)
	}
}

func registerCheck(a domain.App) {
	health.RegisterPeriodicFunc(a.Name, time.Second*time.Duration(a.PollTime), func() error {
		return checkApp(a)
	})
}

func checkApp(a domain.App) error {
	lastApp, _ := getApp(a.ID)
	err := runHTTPCheck(lastApp)
	updateCheckedApp(a, lastApp, err)

	return err
}

func runHTTPCheck(a domain.App) error {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	response, err := client.Get(a.URL)

	if err != nil {
		return err
	}

	if a.CheckType == domain.StatusCheck && a.StatusCode != response.StatusCode {
		return errors.New("Status Code mismatch")
	}

	return nil
}

func updateCheckedApp(a domain.App, lastApp domain.App, err error) {
	var status string
	nowDate := time.Now()

	if err != nil {
		status = "down"
	} else {
		status = "up"
	}

	fmt.Println("App", lastApp.URL, "is", status)

	updateAppStatus(lastApp.ID, status)

	if lastApp.Status != a.Status {
		addHistory(lastApp, nowDate)
		sendNotification(a)
	}
}

func addHistory(app domain.App, date time.Time) {
	if viper.GetBool("history.enabled") {
		history := domain.History{
			AppID:  app.ID,
			Date:   date,
			Status: app.Status,
		}
		insertHistory(history)
	}
}

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

	return err
}
