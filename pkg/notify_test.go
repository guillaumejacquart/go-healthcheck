package pkg

import (
	"testing"

	"github.com/guillaumejacquart/go-healthcheck/pkg/domain"
)

func TestSendNotification(t *testing.T) {
	app := domain.App{
		Name:      "test1",
		URL:       "http://google.fr",
		Status:    "down",
		CheckType: domain.ResponseCheck,
		PollTime:  2,
	}

	err := sendNotification(app)

	if err != nil {
		t.Error("No notification should have been send, so no error should have happened")
	}

	app.Notify = true
	err = sendNotification(app)

	if err == nil {
		t.Error("SMTP error should have happened")
	}
}
