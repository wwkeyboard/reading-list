package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.StaticFS("/css", http.Dir("./static/css"))
	r.StaticFS("/js", http.Dir("./static/js"))
	r.StaticFS("/img", http.Dir("./static/img"))

	r.GET("/health", healthCheck)
	r.GET("/", listReadings)
	r.Run()
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func listReadings(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}