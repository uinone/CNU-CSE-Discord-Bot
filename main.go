package main

import (
	"GO/nomad/utility"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	ds := utility.BotInit()
	defer ds.Close()

	port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

	router := gin.Default()
    router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")

	router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl.html", nil)
    })

	router.Run(":" + port)
}