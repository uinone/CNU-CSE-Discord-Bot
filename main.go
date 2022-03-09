package main

import (
	"GO/nomad/app/controller"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	controller := controller.NewController(time.Minute * 30)

	defer controller.Bot.DiscordSession.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

    if port == "" {
        log.Fatal("$PORT must be set")
    }

	router := gin.Default()
    router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")

	gin.SetMode(gin.ReleaseMode)
	
	router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl.html", nil)
    })

	router.Run(":" + port)
}