package domain

// Header is the header model
type Header struct {
	AppID uint   `json:"app_id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	App   App
}
