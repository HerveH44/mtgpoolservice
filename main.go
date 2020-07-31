package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	database "mtgpoolservice/db"
	"mtgpoolservice/models"
	"mtgpoolservice/models/mtgjson"
	"mtgpoolservice/services"
	"net/http"
)

func init() {
	database.Init()
}

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

		allPrintings := new(mtgjson.AllPrintings)
		if err := json.NewDecoder(resp.Body).Decode(allPrintings); err != nil {
			log.Println("error while unmarshalling allPrintings", err)
		}
		var setsNumber = len(allPrintings.Data)
		log.Println("sets found", setsNumber)

		i := 0
		for setName, set := range allPrintings.Data {
			i++
			fmt.Printf("%d/%d - saving set %s\n", i, setsNumber, setName)
			entity := mtgjson.MapMTGJsonSetToEntity(set)
			if err := database.GetDB().Save(&entity).Error; err != nil {
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

		monoSet := new(mtgjson.MonoSet)
		if err := json.NewDecoder(resp.Body).Decode(monoSet); err != nil {
			log.Println("main: error while unmarshalling monoSet", err)
			return
		}

		log.Println("main: saving set ", monoSet.Data.Name)
		entity := mtgjson.MapMTGJsonSetToEntity(monoSet.Data)
		if err := database.GetDB().Save(&entity).Error; err != nil {
			log.Fatal("main: could not save the set", monoSet.Data.Name, err)

		}
	})
	r.POST("/regular/draft", func(c *gin.Context) {
		var request models.RegularRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ret := make([][]models.Pool, 0)
		for p := 0; p < request.Players; p++ {
			packs, err := services.MakePacks(request.Sets)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			ret = append(ret, packs)
		}
		c.JSON(http.StatusOK, ret)
	})
	r.POST("/regular/sealed", func(c *gin.Context) {
		var request models.RegularRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ret := make([]models.Pool, request.Players)
		for p := 0; p < request.Players; p++ {
			pa := make(models.Pool, 0)
			packs, err := services.MakePacks(request.Sets)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			for _, pack := range packs {
				pa = append(pa, pack...)
			}
			ret[p] = pa
		}
		c.JSON(http.StatusOK, ret)
	})
	r.POST("/cube/sealed", func(c *gin.Context) {
		var request models.CubeSealedRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ret := make([]models.Pool, request.Players)
		for p := 0; p < int(request.Players); p++ {
			pa := make(models.Pool, 0)
			packs, err := services.MakeCubePacks(request)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			for _, pack := range packs {
				pa = append(pa, pack...)
			}
			ret[p] = pa
		}
		c.JSON(http.StatusOK, ret)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
