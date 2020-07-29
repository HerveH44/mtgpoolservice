package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"mtgpoolservice/models"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/refresh", func(context *gin.Context) {
		resp, err := http.Get("http://mtgjson.com/api/v5/AllPrintings.json")
		if err != nil {
			log.Println("Could not fetch the MTGJson allPrintings")
			context.JSON(500, "unexpected error")
			return
		}
		defer resp.Body.Close()

		allPrintings := new(models.AllPrintings)
		if err := json.NewDecoder(resp.Body).Decode(allPrintings); err != nil {
			log.Println("error while unmarshalling allPrintings", err)
		}
		var setsNumber = len(allPrintings.Data)
		log.Println("sets found", setsNumber)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
