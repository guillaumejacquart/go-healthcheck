package pkg

import (
	"testing"
	"time"

	"github.com/guillaumejacquart/go-healthcheck/pkg/domain"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
)

func TestRegisterChecks(t *testing.T) {
	apps, _ := getAllApps()
	for _, a := range apps {
		deleteApp(a.ID)
	}

	app := domain.App{
		Name:      "test1",
		URL:       "http://google.fr",
		Status:    "down",
		CheckType: domain.ResponseCheck,
		PollTime:  2,
	}

	insertApp(&app)

	registerChecks()

	timer := time.NewTimer(time.Second * 3)

	<-timer.C

	app, _ = getApp(app.ID)
	assert.Equal(t, app.Status, "up")
}

func TestRegisterCheck(t *testing.T) {
	app := domain.App{
		Name:      "test2",
		URL:       "http://google.fr",
		Status:    "down",
		CheckType: domain.ResponseCheck,
		PollTime:  2,
	}

	insertApp(&app)

	registerCheck(app)

	timer := time.NewTimer(time.Second * 3)

	<-timer.C

	app, _ = getApp(app.ID)
	assert.Equal(t, app.Status, "up")
}

func TestCheckApp(t *testing.T) {
	app := domain.App{
		Name:      "test3",
		URL:       "http://google.fr",
		Status:    "down",
		CheckType: domain.ResponseCheck,
		PollTime:  5,
	}

	insertApp(&app)

	err := checkApp(app)

	assert.Equal(t, err, nil)

	app, _ = getApp(app.ID)
	assert.Equal(t, app.Status, "up")
}

func TestRunHTTPCheck(t *testing.T) {
	app := domain.App{}
	app.Name = "Test"
	app.URL = "http://google.fr"
	err := runHTTPCheck(app)

	if err != nil {
		t.Error("Google should be checked up")
	}
}

func TestUpdateCheckedApp(t *testing.T) {

	lastApp := domain.App{
		Name:      "test5",
		URL:       "http://google.fr",
		CheckType: domain.ResponseCheck,
		PollTime:  5,
	}

	insertApp(&lastApp)

	app := domain.App{
		Name:      "test6",
		URL:       "http://google.fr",
		CheckType: domain.ResponseCheck,
		PollTime:  5,
	}

	var err error

	updateCheckedApp(app, lastApp, err)

	newApp, newErr := getApp(lastApp.ID)

	if newErr != nil {
		t.Error("Gettin app failed")
	}

	if newApp.Status == "down" {
		t.Error("App status not updated")
	}
}

func TestAddHistory(t *testing.T) {
	viper.Set("history.enabled", true)
	lastApp := domain.App{
		Name:      "test7",
		URL:       "http://google.fr",
		CheckType: domain.ResponseCheck,
		PollTime:  5,
	}

	insertApp(&lastApp)

	addHistory(lastApp, time.Now())

	history, newErr := getAppHistory(lastApp.ID)

	if newErr != nil {
		t.Error("Gettin app failed")
	}

	if len(history) == 0 {
		t.Error("App history not updated")
	}
}

func TestUpdateCheck(t *testing.T) {
	lastApp := domain.App{
		Name:      "test8",
		URL:       "http://google.fr",
		CheckType: domain.ResponseCheck,
		PollTime:  5,
	}

	insertApp(&lastApp)

	lastApp.CheckStatus = "start"
	updateCheck(lastApp)

	timer := timers[lastApp.ID]
	if timer == nil {
		t.Error("Timer is nul after app status set to start")
	}

	lastApp.CheckStatus = "stop"
	updateCheck(lastApp)

	_, ok := timers[lastApp.ID]

	assert.Equal(t, ok, false)
}
