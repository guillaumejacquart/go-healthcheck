package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func runChecksApp(c chan App) {
	apps, err := getAllApps()

	if err != nil {
		panic(err)
	}

	for _, a := range apps {
		go func(a App) {
			runHTTPCheck(a, c)
		}(a)
	}

	for a := range c {
		go runHTTPCheck(a, c)
	}
}

func runHTTPCheck(a App, c chan App) {
	nowDate := time.Now()
	lastApp, _ := getApp(a.ID)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	_, err := client.Get(lastApp.URL)

	if err != nil {
		lastApp.Status = "down"
	} else {
		lastApp.Status = "up"
		lastApp.LastUpDate = nowDate
	}

	fmt.Println("App", lastApp.URL, "is", lastApp.Status)

	updateApp(lastApp.ID, lastApp)

	if lastApp.Status != a.Status {
		addHistory(lastApp, nowDate)
	}

	time.Sleep(time.Second * time.Duration(lastApp.PollTime))
	c <- lastApp
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
