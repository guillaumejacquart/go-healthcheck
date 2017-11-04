package main

import "time"

// App is the app model
type App struct {
	Name       string    `json:"name" binding:"required"`
	URL        string    `json:"url"`
	Status     string    `json:"status"`
	PollTime   int       `json:"pollTime"`
	LastUpDate time.Time `json:"lastUpDate"`
}
