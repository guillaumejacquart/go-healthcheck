package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var appChan chan App

func initApi(c chan App) {
	appChan = c
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Next()
	}
}

// Serve api server to specified port
func Serve(port int) {
	router := gin.Default()

	router.Use(cors.Default())
	router.Use(Logger())

	router.Static("/app", "./public")

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/apps", func(c *gin.Context) {
		apps := getAllApps()
		c.JSON(200, apps)
	})

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/apps/:name", func(c *gin.Context) {
		name := c.Param("name")
		app, err := getApp(name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, app)
	})

	router.POST("/apps", func(c *gin.Context) {
		var app App
		if err := c.BindJSON(&app); err == nil {
			fmt.Println(app)
			err := insertApp(app)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				go runHttpCheck(app, appChan)
				c.JSON(http.StatusOK, app)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	router.PUT("/apps/:name", func(c *gin.Context) {
		name := c.Param("name")
		var app App
		if err := c.BindJSON(&app); err == nil {
			fmt.Println(app)
			err := updateApp(name, app)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, app)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	router.DELETE("/apps/:name", func(c *gin.Context) {
		name := c.Param("name")
		err := deleteApp(name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	router.Run(fmt.Sprintf(":%v", port))
}
