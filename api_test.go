package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestApiGetAllApps(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/apps", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
}

func TestApiGetAppNotExist(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/apps/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusInternalServerError)
}

func TestApiGetAppExist(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	app := App{
		Name:      "test",
		URL:       "http://google.fr",
		CheckType: responseCheck,
	}

	insertApp(&app)

	req, _ := http.NewRequest("GET", "/apps/"+fmt.Sprint(app.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	exApp := new(App)
	err := json.Unmarshal(w.Body.Bytes(), &exApp)

	assert.Equal(t, err, nil)

	assert.Equal(t, exApp.Name, app.Name)
}

func TestApiCreateApp(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	app := App{
		Name:      "test",
		URL:       "http://google.fr",
		CheckType: responseCheck,
		PollTime:  5,
	}

	appBytes, _ := json.Marshal(app)

	req, _ := http.NewRequest("POST", "/apps", bytes.NewReader(appBytes))
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
	exApps, err := getAllApps()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, exApps[0].Name, app.Name)
}

func TestApiUpdateApp(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	app := App{
		Name:      "test",
		URL:       "http://google.fr",
		CheckType: responseCheck,
		PollTime:  5,
	}

	insertApp(&app)

	app.URL = "http://amazon.fr"

	appBytes, _ := json.Marshal(app)

	req, _ := http.NewRequest("PUT", "/apps/"+fmt.Sprint(app.ID), bytes.NewReader(appBytes))
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
	exApp, err := getApp(app.ID)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, exApp.URL, app.URL)
}

func TestApiCreateAppHistory(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	app := App{
		Name:      "test",
		URL:       "http://google.fr",
		CheckType: responseCheck,
		PollTime:  5,
	}

	insertApp(&app)

	history := History{
		AppID:  app.ID,
		Date:   time.Now(),
		Status: "up",
	}

	insertHistory(history)

	req, _ := http.NewRequest("GET", "/apps/"+fmt.Sprint(app.ID)+"/history", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	histories := []History{}
	err := json.Unmarshal(w.Body.Bytes(), &histories)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(histories), 1)

	assert.Equal(t, histories[0].AppID, app.ID)

	if err != nil {
		t.Error(err)
	}
}

func TestApiDeleteApp(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()

	app := App{
		Name:      "test",
		URL:       "http://google.fr",
		CheckType: responseCheck,
	}

	insertApp(&app)

	req, _ := http.NewRequest("DELETE", "/apps/"+fmt.Sprint(app.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	_, err := getApp(app.ID)

	if err == nil {
		t.Error("App deletion failed")
	}
}
