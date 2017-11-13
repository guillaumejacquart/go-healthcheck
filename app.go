package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type checkType int

const (
	responseCheck = checkType(0)
	statusCheck   = checkType(1)
)

// App is the app model
type App struct {
	gorm.Model
	Name       string    `json:"name" binding:"required"`
	URL        string    `json:"url"`
	Status     string    `json:"status"`
	PollTime   int       `json:"pollTime"`
	LastUpDate time.Time `json:"lastUpDate"`
	CheckType  checkType `json:"checkType"`
	StatusCode int       `json:"statusCode"`
	Notify     bool      `json:"notify"`
}
