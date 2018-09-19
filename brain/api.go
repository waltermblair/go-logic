package brain

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RunAPI(rabbit RabbitClient) {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/run", func(c *gin.Context) {
		var msg Message

		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		body := msg.Body

		c.JSON(http.StatusOK, gin.H{"config": body.Configs, "input": body.Input})

		if err := RunDemo(body, rabbit); err != nil {
			fmt.Println("error running app: ", err)
		}
	})

	r.Use(cors.Default())
	r.Run() // listen and serve on 0.0.0.0:8080

}