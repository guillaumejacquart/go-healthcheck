package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func startCheck() {
	c := make(chan App)
	initApi(c)
	initDb("db")
	go runChecksApp(c)
	Serve(8080)
}

func runChecksApp(c chan App) {
	apps := getAllApps()

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
	lastApp, _ := getApp(a.Name)
	_, err := http.Get(lastApp.URL)

	if err != nil {
		lastApp.Status = "down"
	} else {
		lastApp.Status = "up"
		lastApp.LastUpDate = nowDate
	}

	fmt.Println("App", lastApp.URL, "is", lastApp.Status)

	updateApp(lastApp.Name, lastApp)

	if lastApp.Status != a.Status {
		insertHistory(lastApp, nowDate)
	}

	time.Sleep(time.Second * time.Duration(lastApp.PollTime))
	c <- lastApp
}

func insertHistory(app App, date time.Time) {
	if viper.GetBool("history.enabled") {
		history := History{
			AppName: app.Name,
			Date:    date,
			Status:  app.Status,
		}
		dbType := viper.Get("history.db_type")
		if dbType == "file" {
			insertObject("history", fmt.Sprintf("%v", date.Unix()), history)
		}
	}
}
