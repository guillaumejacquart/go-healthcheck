package main

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"

	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	initConfig()

	viper.Set("db.type", "sqlite3")
	viper.Set("db.path", "test.db")

	initDb()

	retCode := m.Run()

	log.Print("Dropping tables...")
	db.DropTable("apps")
	db.DropTable("histories")
	log.Print("Tables dropped !")

	os.Exit(retCode)
}

func TestRunChecksApp(t *testing.T) {
	apps, _ := getAllApps()
	for _, a := range apps {
		deleteApp(a.ID)
	}

	app := App{
		Name:      "test1",
		URL:       "http://google.fr",
		Status:    "down",
		CheckType: responseCheck,
		PollTime:  2,
	}

	insertApp(&app)

	runChecksApp()

	timer := time.NewTimer(time.Second * 3)

	<-timer.C

	app, _ = getApp(app.ID)
	assert.Equal(t, app.Status, "up")
}

func TestRegisterCheck(t *testing.T) {
	app := App{
		Name:      "test2",
		URL:       "http://google.fr",
		Status:    "down",
		CheckType: responseCheck,
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
	app := App{
		Name:      "test3",
		URL:       "http://google.fr",
		Status:    "down",
		CheckType: responseCheck,
		PollTime:  5,
	}

	insertApp(&app)

	err := checkApp(app)

	assert.Equal(t, err, nil)

	app, _ = getApp(app.ID)
	assert.Equal(t, app.Status, "up")
}

func TestRunHTTPCheck(t *testing.T) {
	app := App{}
	app.Name = "Test"
	app.URL = "http://google.fr"
	err := runHTTPCheck(app)

	if err != nil {
		t.Error("Google should be checked up")
	}
}

func TestUpdateCheckedApp(t *testing.T) {

	lastApp := App{
		Name:      "test5",
		URL:       "http://google.fr",
		CheckType: responseCheck,
		PollTime:  5,
	}

	insertApp(&lastApp)

	app := App{
		Name:      "test6",
		URL:       "http://google.fr",
		CheckType: responseCheck,
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
	lastApp := App{
		Name:      "test7",
		URL:       "http://google.fr",
		CheckType: responseCheck,
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
