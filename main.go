package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mtgpoolservice/common"
	"mtgpoolservice/models"
	"mtgpoolservice/services"
	"net/http"
)

var db = common.Init()

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/refresh", func(context *gin.Context) {
		log.Println("Refresh data")
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

		i := 0
		for setName, set := range allPrintings.Data {
			i++
			fmt.Printf("%d/%d - saving set %s\n", i, setsNumber, setName)
			if err := db.Save(&set).Error; err != nil {
				fmt.Printf("could not save the card %s - %s\n", setName, err)
			}
		}
	})
	r.GET("/test", func(context *gin.Context) {
		log.Println("Refresh data")
		resp, err := http.Get("http://mtgjson.com/api/v5/ISD.json")
		if err != nil {
			log.Println("Could not fetch the MTGJson monoSet")
			context.JSON(500, "unexpected error: Could not fetch the MTGJson monoSet")
			return
		}
		defer resp.Body.Close()

		monoSet := new(models.MonoSet)
		if err := json.NewDecoder(resp.Body).Decode(monoSet); err != nil {
			log.Println("main: error while unmarshalling monoSet", err)
		}

		log.Println("saving cards from ", monoSet.Data.Name)
		if err := db.Save(&monoSet.Data).Error; err != nil {
			log.Fatal("main: could not save the card", monoSet.Data.Name, err)
		}
	})
	r.POST("/regular/draft", func(c *gin.Context) {
		var request models.RegularDraftRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ret := make([][]models.Pack, 0)
		for p := 0; p < request.Players; p++ {
			packs, err := services.MakePacks(request.Sets)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
			}
			ret = append(ret, packs)
		}
		c.JSON(http.StatusOK, ret)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
