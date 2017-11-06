package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	health "github.com/docker/go-healthcheck"
	"github.com/spf13/viper"
)

func runChecksApp(c chan App) {
	apps, err := getAllApps()

	if err != nil {
		panic(err)
	}

	for _, a := range apps {
		registerCheck(a)
	}
}

func registerCheck(a App) {
	health.RegisterPeriodicFunc(a.Name, time.Second*time.Duration(a.PollTime), func() error {
		lastApp, _ := getApp(a.ID)
		err := runHTTPCheck(lastApp)
		updateCheckedApp(a, lastApp, err)

		return err
	})
}

func runHTTPCheck(a App) error {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	response, err := client.Get(a.URL)

	if err != nil {
		return err
	}

	if a.CheckType == statusCheck && a.StatusCode != response.StatusCode {
		return errors.New("Status Code mismatch")
	}

	return nil
}

func updateCheckedApp(a App, lastApp App, err error) {
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
	}
}

func addHistory(app App, date time.Time) {
	if viper.GetBool("history.enabled") {
		history := History{
			AppID:  app.ID,
			Date:   date,
			Status: app.Status,
		}
		insertHistory(history)
	}
}
