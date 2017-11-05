package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var db *gorm.DB

func initDb() {
	dbType := viper.GetString("db.type")
	dbPort := viper.GetInt("db.port")
	dbHost := viper.GetString("db.host")
	dbUsername := viper.GetString("db.username")
	dbPassword := viper.GetString("db.password")
	dbName := viper.GetString("db.name")

	var connectionString string
	switch dbType {
	case "mysql":
		connectionString = fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/%v?parseTime=true",
			dbUsername,
			dbPassword,
			dbHost,
			dbPort,
			dbName)
	case "postgres":
		connectionString = fmt.Sprintf(
			"host=%v port=%v user=%v dbname=%v sslmode=disable password=%v",
			dbHost,
			dbPort,
			dbUsername,
			dbName,
			dbPassword)
	}

	dbInit, err := gorm.Open(dbType, connectionString)
	if err != nil {
		panic(err)
	}

	db = dbInit

	// Migrate the schema
	db.AutoMigrate(&App{}, &History{})
}

func getAllApps() ([]App, error) {
	var apps []App
	err := db.Find(&apps).Error

	return apps, err
}

func getApp(id uint) (App, error) {
	app := App{}
	err := db.First(&app, id).Error

	return app, err
}

func insertApp(app *App) error {
	app.LastUpDate = time.Now()
	return db.Create(app).Error
}

func insertHistory(history History) error {
	return db.Create(&history).Error
}

func getAppHistory(appID uint) ([]History, error) {
	histories := []History{}
	err := db.Order("date desc").Limit(5).Where("app_id = ?", appID).Find(&histories).Error

	return histories, err
}

func updateApp(id uint, app App) error {
	existingApp, err := getApp(id)
	if err != nil {
		return err
	}

	existingApp.URL = app.URL
	existingApp.PollTime = app.PollTime

	return db.Save(&existingApp).Error
}

func updateAppStatus(id uint, status string) error {
	existingApp, err := getApp(id)
	if err != nil {
		return err
	}

	existingApp.Status = status
	if status == "up" {
		existingApp.LastUpDate = time.Now()
	}

	return db.Save(&existingApp).Error
}

func deleteApp(id uint) error {
	app, err := getApp(id)
	if err != nil {
		return err
	}

	return db.Delete(&app).Error
}
