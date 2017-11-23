package domain

import (
	"time"

	"github.com/jinzhu/gorm"
)

type checkType int

const (
	ResponseCheck = checkType(0)
	StatusCheck   = checkType(1)
)

// App is the app model
type App struct {
	gorm.Model
	Name        string     `json:"name" binding:"required"`
	URL         string     `json:"url"`
	Status      string     `json:"status"`
	CheckStatus string     `json:"checkStatus"`
	PollTime    int        `json:"pollTime"`
	LastUpDate  *time.Time `json:"lastUpDate"`
	CheckType   checkType  `json:"checkType"`
	StatusCode  int        `json:"statusCode"`
	Notify      bool       `json:"notify"`
	Headers     []Header   `json:"headers"`
}
