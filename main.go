package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wwkeyboard/reading-list/reading"
)

var (
	db *reading.Database
)

func init() {
	newDB, err := reading.NewDatabase("reading.sql")
	db = newDB
	if err != nil {
		log.Panicf("Couldn't load DB %+v", err)
	}

	err = db.EnsureBucket()
	if err != nil {
		log.Panicf("Couldn't create default bucket %+v", err)
	}

}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.StaticFS("/css", http.Dir("./static/css"))
	r.StaticFS("/js", http.Dir("./static/js"))
	r.StaticFS("/img", http.Dir("./static/img"))

	r.GET("/health", healthCheck)
	r.GET("/", listReadings)
	r.POST("/", createPiece)

	r.Run()
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func listReadings(c *gin.Context) {
	ps, _ := db.List()

	c.JSON(http.StatusOK, gin.H{"pieces": ps})
}

func createPiece(c *gin.Context) {
	name := c.PostForm("name")
	url := c.PostForm("url")

	db.AddPiece(&reading.Piece{
		Name: name,
		URL:  url,
	})
}
