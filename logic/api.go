package logic

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RunAPI() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Use(cors.Default())
	r.Run("0.0.0.0:8080") // listen and serve on 0.0.0.0:8080

}