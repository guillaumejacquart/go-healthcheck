package api

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Remove("_test")
	initDb("_test")
	retCode := m.Run()

	os.Remove("_test")
	os.Exit(retCode)
}

func TestInitDb(t *testing.T) {
	_, err := os.Open("_test")
	if err != nil {
		t.Error(err)
	}
}

func TestInsertApp(t *testing.T) {
	insertApp(App{Name: "_test"})

	_, err := os.Open("_test/apps/_test.json")

	if err != nil {
		t.Error("File for app does not exist")
	}

	apps := getAllApps()

	if len(apps) <= 0 {
		t.Error("App not created successfully")
	}
}
