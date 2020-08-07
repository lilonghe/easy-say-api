package main

import (
	"github.com/easy-say-api/src/config"
	"github.com/easy-say-api/src/router/message"
	"github.com/easy-say-api/src/router/user"
	ginsession "github.com/go-session/gin-session"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(ginsession.New())

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/login", user.Login)

	authorized := r.Group("/")
	authorized.Use(AuthRequired())
	{
		authorized.GET("/message/:id", message.Get)
		authorized.GET("/message", message.GetList)
		authorized.POST("/message", message.Add)
	}

	return r
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		store := ginsession.FromContext(c)
		openid, _ := store.Get("openid")
		if openid == nil {
			c.JSON(200, gin.H{"err": "未登录"})
			c.Abort()
			return
		}
	}
}

func main() {
	config.Init()

	r := setupRouter()
	r.Run(":8080")
}
