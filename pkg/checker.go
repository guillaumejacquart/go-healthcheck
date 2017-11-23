package pkg

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/guillaumejacquart/go-healthcheck/pkg/domain"
	"github.com/spf13/viper"
)

var timers map[uint]*time.Ticker = make(map[uint]*time.Ticker)

func registerChecks() {
	apps, err := getAllApps()

	if err != nil {
		panic(err)
	}

	for _, a := range apps {
		if a.CheckStatus == "start" {
			registerCheck(a)
		}
	}
}

func registerCheck(a domain.App) {
	ticker := time.NewTicker(time.Second * time.Duration(a.PollTime))
	go func() {
		for range ticker.C {
			checkApp(a)
		}
	}()
	timers[a.ID] = ticker
}

func updateCheck(a domain.App) {
	switch a.CheckStatus {
	case "stop":
		timers[a.ID].Stop()
		delete(timers, a.ID)
	case "start":
		registerCheck(a)
	}
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
	req, err := http.NewRequest("GET", a.URL, nil)

	if len(a.Headers) > 0 {
		for _, h := range a.Headers {
			req.Header.Add(h.Name, h.Value)
		}
	}

	response, err := client.Do(req)

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
