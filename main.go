package main

import (
	"github.com/easy-say-api/src/config"
	"github.com/easy-say-api/src/router/message"
	"github.com/easy-say-api/src/router/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginsession "github.com/go-session/gin-session"
	"time"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"OPTION", "PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
			//return origin == "http://127.0.0.1"
		},
		MaxAge: 12 * time.Hour,
	}))
	r.Use(ginsession.New())

	r.POST("/api/login", user.Login)

	authorized := r.Group("/api")
	authorized.Use(AuthRequired())
	{
		authorized.GET("/user", user.UserInfo)
		authorized.GET("/message/:id", message.Get)
		authorized.GET("/message", message.GetList)
		authorized.POST("/message", message.Add)
		authorized.POST("/message/like", message.Like)
	}

	return r
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		store := ginsession.FromContext(c)
		openid, _ := store.Get("openid")
		c.Set("openid", openid)
		if openid == nil {
			c.JSON(401, gin.H{"err": "未登录"})
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
