package main

import "time"

// History is the history model
type History struct {
	AppName string    `json:"app_name" binding:"required"`
	Status  string    `json:"status"`
	Date    time.Time `json:"date"`
}
