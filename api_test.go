package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magiconair/properties/assert"
)

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
