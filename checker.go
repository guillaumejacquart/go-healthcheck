package main

import (
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
		return runHTTPCheck(a)
	})
}

func runHTTPCheck(a App) error {
	nowDate := time.Now()
	lastApp, _ := getApp(a.ID)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	_, err := client.Get(lastApp.URL)

	var status string

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

	return err
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
