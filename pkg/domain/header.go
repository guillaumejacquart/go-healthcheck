package domain

import (
	"github.com/jinzhu/gorm"
)

// Header is the header model
type Header struct {
	gorm.Model
	AppID uint   `json:"app_id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
