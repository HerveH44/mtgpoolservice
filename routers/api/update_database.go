package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	database "mtgpoolservice/db"
	"mtgpoolservice/models/mtgjson"
	"net/http"
)

func Refresh(context *gin.Context) {
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
}

func RefreshSet(context *gin.Context) {
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
}
