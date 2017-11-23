package pkg

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/guillaumejacquart/go-healthcheck/pkg/domain"
)

type Server struct {
	Router *gin.Engine
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Next()
	}
}

func createServer() Server {
	server := Server{
		Router: gin.Default(),
	}
	return server
}

func (s *Server) initializeMiddlewares() {
	s.Router.Use(cors())
	s.Router.Use(gin.Recovery())
}

// Serve api server to specified port
func (s *Server) setupRoutes() {
	router := s.Router
	router.Static("/app", "./public")

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/apps", func(c *gin.Context) {
		apps, err := getAllApps()

		if err != nil {
			panic(err)
		}

		c.JSON(200, apps)
	})

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/apps/:id", func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)

		if err != nil {
			panic(err)
		}

		app, err := getApp(uint(idInt))
		if err != nil {
			panic(err)
		}
		c.JSON(200, app)
	})

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/apps/:id/history", func(c *gin.Context) {
		id := getId(c)

		histories, err := getAppHistory(uint(id))

		if err != nil {
			panic(err)
		}

		c.JSON(200, histories)
	})

	router.POST("/apps", func(c *gin.Context) {
		var app = domain.App{}
		if err := c.BindJSON(&app); err != nil {
			panic(err)
		}

		err := insertApp(&app)
		if err != nil {
			panic(err)
		} else {
			registerCheck(app)
			c.JSON(http.StatusOK, app)
		}
	})

	router.PUT("/apps/:id", func(c *gin.Context) {
		id := getId(c)

		var app domain.App
		if err := c.BindJSON(&app); err != nil {
			panic(err)
		}

		oldApp, _ := getApp(uint(id))

		err := updateApp(uint(id), app)
		if err != nil {
			panic(err)
		} else {
			if app.CheckStatus != oldApp.CheckStatus {
				lastApp, _ := getApp(uint(id))
				updateCheck(lastApp)
			}
			c.JSON(http.StatusOK, app)
		}
	})

	router.DELETE("/apps/:id", func(c *gin.Context) {
		id := getId(c)

		err := deleteApp(uint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})
}

func getId(c *gin.Context) int {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)

	if err != nil {
		panic(err)
	}

	return idInt
}

func (s *Server) serve(port int) {
	s.Router.Run(fmt.Sprintf(":%v", port))
}
