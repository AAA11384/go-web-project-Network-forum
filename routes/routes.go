package routes

import (
	"net/http"
	"qimiproject/controllers"
	"qimiproject/logger"
	"qimiproject/middlewares"
	"time"

	"github.com/gin-gonic/gin"
)

func SetUp() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	v1 := r.Group("/api/v1")
	v1.POST("/signup", controllers.SignUpHandler)

	v1.POST("/login", controllers.LoginHandler)

	v1.Use(middlewares.BucketLimit(1, time.Microsecond*500), middlewares.JWTAuthMiddleware())

	{
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)

		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.PostDetailHandler)
		v1.GET("/posts", controllers.PostListHandler)
		v1.GET("/community_posts", controllers.CommunityPostsOrderHandler)

		v1.POST("/vote", controllers.PostVoteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": 404,
		})
	})
	return r
}
