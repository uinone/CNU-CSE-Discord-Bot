package model

import (
	"GO/nomad/app/view"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/bwmarrin/discordgo"
)

type Web struct {
	port  	string
	router 	*gin.Engine
	viwer 	*view.Viewer
}

// Create web object
func NewWeb(ds *discordgo.Session) *Web {
	w := new(Web)

	w.viwer = view.NewViewer()
	w.viwer.SetDiscordSession(ds)

	w.port = os.Getenv("PORT")
	if w.port == "" {
		w.port = "3000"
	}	

    if w.port == "" {
        w.viwer.FatallnMsgToConsole("$PORT must be set")
    }

	w.router = gin.Default()
	w.router.Use(gin.Logger())
	w.router.LoadHTMLGlob("templates/*.tmpl.html")

	return w
}

// Run router of gin
func (w *Web) Run() {
	gin.SetMode(gin.ReleaseMode)

	w.router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl.html", nil)
    })

	w.router.Run(":" + w.port)
}