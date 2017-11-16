package pkg

import (
	"testing"

	"github.com/guillaumejacquart/go-healthcheck/pkg/domain"
)

func insertTestApp() (*domain.App, error) {

	testApp := new(domain.App)
	testApp.Name = "Test"
	testApp.URL = "http://google.fr"
	testApp.Status = "down"

	err := insertApp(testApp)
	return testApp, err
}

func truncateTestTables() {
	db.Exec("delete from apps")
	db.Exec("delete from histories")
}

func TestInitDb(t *testing.T) {
	if db == nil {
		t.Error("Init db failed")
	}
}

func TestInsertApp(t *testing.T) {
	app, err := insertTestApp()
	if err != nil {
		t.Error("Error occured on insertion", err)
	}

	if app.ID == 0 {
		t.Error("App ID not created after insert")
	}
	truncateTestTables()
}

func TestGetAllApps(t *testing.T) {
	insertTestApp()

	apps, err := getAllApps()

	if err != nil {
		t.Error("App get failed")
	}

	if len(apps) != 1 {
		t.Error("App table length incorrect")
	}
	truncateTestTables()
}

func TestGetApp(t *testing.T) {
	app, _ := insertTestApp()

	getApp, err := getApp(app.ID)

	if err != nil {
		t.Error("App get failed")
	}

	if getApp.Name != "Test" {
		t.Error("App name incorrect")
	}
	truncateTestTables()
}

func TestDeleteApp(t *testing.T) {
	app, _ := insertTestApp()

	err := deleteApp(app.ID)

	if err != nil {
		t.Error("App deletion failed")
	}

	_, err = getApp(app.ID)

	if err == nil {
		t.Error("App still exists")
	}

	truncateTestTables()
}
