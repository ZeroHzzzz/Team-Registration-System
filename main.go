package main

import (
	"backend/config/database"
	"backend/config/router"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Init()

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	router.Init(r)
	err := r.Run(":3000")
	if err != nil {
		log.Fatal("Server start error", err)
	}
}
