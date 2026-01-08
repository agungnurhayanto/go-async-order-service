package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// Middleware dasar

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// healt check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})

	})

	log.Println("Api running on :8090")

	if err := r.Run(":8090"); err != nil {
		log.Fatal(err)
	}
}
