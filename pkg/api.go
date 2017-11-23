package pkg

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/guillaumejacquart/go-healthcheck/pkg/domain"
	"github.com/spf13/viper"
)

type Server struct {
	Router    *gin.Engine
	ApiRouter *gin.RouterGroup
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
	s.ApiRouter = s.Router.Group("/api")

	if viper.GetBool("authentication.enabled") {
		// the jwt middleware
		authMiddleware := &jwt.GinJWTMiddleware{
			Realm:      "test zone",
			Key:        []byte("secret key"),
			Timeout:    time.Hour,
			MaxRefresh: time.Hour,
			Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
				if userId == viper.GetString("authentication.username") && password == viper.GetString("authentication.password") {
					return userId, true
				}

				return userId, false
			},
			Authorizator: func(userId string, c *gin.Context) bool {
				return true
			},
			Unauthorized: func(c *gin.Context, code int, message string) {
				c.JSON(code, gin.H{
					"code":    code,
					"message": message,
				})
			},
			TokenLookup:   "header:Authorization",
			TokenHeadName: "Bearer",
			TimeFunc:      time.Now,
		}
		s.Router.POST("/login", authMiddleware.LoginHandler)
		s.ApiRouter.Use(authMiddleware.MiddlewareFunc())
	}
}

// Serve api server to specified port
func (s *Server) setupRoutes() {
	s.Router.Static("/app", "./public")
	router := s.ApiRouter

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
		id := getId(c.Param("id"))

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
		id := getId(c.Param("id"))

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
		id := getId(c.Param("id"))

		err := deleteApp(uint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})
}

func getId(id string) int {
	idInt, err := strconv.Atoi(id)

	if err != nil {
		panic(err)
	}

	return idInt
}

func (s *Server) serve(port int) {
	s.Router.Run(fmt.Sprintf(":%v", port))
}
