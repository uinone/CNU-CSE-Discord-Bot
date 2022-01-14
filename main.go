package main

import (
	"GO/nomad/utility"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
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

	ds := utility.BotInit()
	defer ds.Close()

	utility.RunAlarm(ds, time.Minute * 30)
}