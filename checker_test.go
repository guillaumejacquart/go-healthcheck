package main

import (
	"log"
	"os"
	"testing"
	"time"

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

	lastApp := new(App)
	lastApp.Name = "Test"
	lastApp.URL = "http://google.fr"
	lastApp.Status = "down"

	insertApp(lastApp)

	app := App{}
	app.Name = "Test"
	app.URL = "http://google.fr"
	app.Status = "up"

	var err error

	updateCheckedApp(app, *lastApp, err)

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
	lastApp := new(App)
	lastApp.Name = "Test"
	lastApp.URL = "http://google.fr"
	lastApp.Status = "down"

	insertApp(lastApp)

	addHistory(*lastApp, time.Now())

	history, newErr := getAppHistory(lastApp.ID)

	if newErr != nil {
		t.Error("Gettin app failed")
	}

	if len(history) == 0 {
		t.Error("App history not updated")
	}
}
