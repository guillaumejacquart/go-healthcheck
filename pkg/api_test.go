package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guillaumejacquart/go-healthcheck/pkg/domain"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
)

func getRouter() *gin.Engine {
	server := createServer()
	server.initializeMiddlewares()
	server.setupRoutes()

	return server.Router
}

func TestApiGetAllApps(t *testing.T) {
	router := getRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/apps", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
}

func TestGetIdPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	getId("notint")
}

func TestApiGetAppNotExist(t *testing.T) {
	router := getRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/apps/1234", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusInternalServerError)
}

func TestApiGetAppExist(t *testing.T) {
	router := getRouter()

	w := httptest.NewRecorder()

	app := domain.App{
		Name:      "test",
		URL:       "http://google.fr",
		CheckType: domain.ResponseCheck,
	}

	insertApp(&app)

	req, _ := http.NewRequest("GET", "/api/apps/"+fmt.Sprint(app.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	exApp := new(domain.App)
	err := json.Unmarshal(w.Body.Bytes(), &exApp)

	assert.Equal(t, err, nil)

	assert.Equal(t, exApp.Name, app.Name)
}

func TestApiCreateApp(t *testing.T) {
	router := getRouter()

	w := httptest.NewRecorder()

	app := domain.App{
		Name:      "test",
		URL:       "http://google.fr",
		CheckType: domain.ResponseCheck,
		PollTime:  5,
	}

	appBytes, _ := json.Marshal(app)

	req, _ := http.NewRequest("POST", "/api/apps", bytes.NewReader(appBytes))
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
	exApps, err := getAllApps()

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, exApps[0].Name, app.Name)
}

func TestApiUpdateApp(t *testing.T) {
	router := getRouter()

	w := httptest.NewRecorder()

	app := domain.App{
		Name:      "test",
		URL:       "http://google.fr",
		CheckType: domain.ResponseCheck,
		PollTime:  5,
	}

	insertApp(&app)

	app.URL = "http://amazon.fr"

	appBytes, _ := json.Marshal(app)

	req, _ := http.NewRequest("PUT", "/api/apps/"+fmt.Sprint(app.ID), bytes.NewReader(appBytes))
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)
	exApp, err := getApp(app.ID)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, exApp.URL, app.URL)
}

func TestApiCreateAppHistory(t *testing.T) {
	router := getRouter()

	w := httptest.NewRecorder()

	app := domain.App{
		Name:      "test",
		URL:       "http://google.fr",
		CheckType: domain.ResponseCheck,
		PollTime:  5,
	}

	insertApp(&app)

	history := domain.History{
		AppID:  app.ID,
		Date:   time.Now(),
		Status: "up",
	}

	insertHistory(history)

	req, _ := http.NewRequest("GET", "/api/apps/"+fmt.Sprint(app.ID)+"/history", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	histories := []domain.History{}
	err := json.Unmarshal(w.Body.Bytes(), &histories)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(histories), 1)

	assert.Equal(t, histories[0].AppID, app.ID)

	if err != nil {
		t.Error(err)
	}
}

func TestApiDeleteApp(t *testing.T) {
	router := getRouter()

	w := httptest.NewRecorder()

	app := domain.App{
		Name:      "test",
		URL:       "http://google.fr",
		CheckType: domain.ResponseCheck,
	}

	insertApp(&app)

	req, _ := http.NewRequest("DELETE", "/api/apps/"+fmt.Sprint(app.ID), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

	_, err := getApp(app.ID)

	if err == nil {
		t.Error("App deletion failed")
	}
}

func TestApiAuthorization(t *testing.T) {
	viper.Set("authentication.enabled", true)
	viper.Set("authentication.username", "admin")
	viper.Set("authentication.password", "admin")
	router := getRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/api/apps", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusUnauthorized)

	viper.Set("authentication.enabled", false)
}
