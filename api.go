package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var appChan chan App

func initAPI(c chan App) {
	appChan = c
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Next()
	}
}

// Serve api server to specified port
func Serve(port int) {
	router := gin.Default()

	router.Use(cors())

	router.Static("/app", "./public")

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/apps", func(c *gin.Context) {
		apps, err := getAllApps()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, apps)
	})

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/apps/:id", func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		app, err := getApp(uint(idInt))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, app)
	})

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/apps/:id/history", func(c *gin.Context) {
		id := c.Param("id")
		idInt, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		histories, err := getAppHistory(uint(idInt))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, histories)
	})

	router.POST("/apps", func(c *gin.Context) {
		var app App
		if err := c.BindJSON(&app); err == nil {
			fmt.Println(app)
			err := insertApp(&app)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				go runHTTPCheck(app, appChan)
				c.JSON(http.StatusOK, app)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	router.PUT("/apps/:id", func(c *gin.Context) {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var app App
		if err := c.BindJSON(&app); err == nil {
			fmt.Println(app)
			err := updateApp(uint(idInt), app)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, app)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	router.DELETE("/apps/:id", func(c *gin.Context) {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = deleteApp(uint(idInt))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	router.Run(fmt.Sprintf(":%v", port))
}
