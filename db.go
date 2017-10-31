package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nanobox-io/golang-scribble"
)

var db *scribble.Driver

func initDb(path string) {
	newDb, err := scribble.New(path, nil)
	db = newDb

	if err != nil {
		fmt.Println("Error", err)
	}
}

func getAllApps() []App {
	records, err := db.ReadAll("apps")
	if err != nil {
		fmt.Println("Error", err)
	}

	apps := []App{}
	for _, f := range records {
		appFound := App{}
		if err := json.Unmarshal([]byte(f), &appFound); err != nil {
			fmt.Println("Error", err)
		}
		apps = append(apps, appFound)
	}

	return apps
}

func getApp(name string) (App, error) {
	app := App{}
	if err := db.Read("apps", name, &app); err != nil {
		return app, err
	}

	return app, nil
}

func insertApp(app App) error {
	existingApp, err := getApp(app.Name)
	if err != nil {
		fmt.Println(err)
	}
	if existingApp.Name == app.Name {
		return errors.New("app name already exists")
	}
	return db.Write("apps", app.Name, app)
}

func updateApp(name string, app App) error {
	_, err := getApp(name)
	if err != nil {
		return err
	}
	return db.Write("apps", name, app)
}

func deleteApp(name string) error {
	_, err := getApp(name)
	if err != nil {
		return err
	}
	return db.Delete("apps", name)
}
