package domain

import "time"

// History is the history model
type History struct {
	AppID  uint      `json:"app_id"`
	Status string    `json:"status"`
	Date   time.Time `json:"date"`
}
