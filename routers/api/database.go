package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	database "mtgpoolservice/db"
	"mtgpoolservice/models"
	"mtgpoolservice/models/entities"
	"mtgpoolservice/models/mtgjson"
	"net/http"
	"sort"
)

var playableSetTypes = []string{"core", "expansion", "draft_innovation", "funny", "starter", "masters"}

func AvailableSets(context *gin.Context) {
	sets, err := database.GetSets()
	if err != nil {
		context.JSON(500, gin.H{"error": "unexpected error"})
		return
	}

	setMap := buildAvailableSetsMap(sets[:])
	context.JSON(http.StatusOK, setMap)
}

func buildAvailableSetsMap(sets []entities.Set) map[string][]models.SetResponse {
	setMap := setupSetsMap()

	for _, set := range sets {
		setType := set.Type
		i := sort.SearchStrings(playableSetTypes, setType)
		if i >= len(playableSetTypes) || playableSetTypes[i] != setType {
			continue
		}
		setMap[setType] = append(setMap[setType], models.SetResponse{
			Code: set.Code,
			Name: set.Name,
		})
	}
	return setMap
}

func setupSetsMap() map[string][]models.SetResponse {
	setMap := make(map[string][]models.SetResponse)
	for _, t := range playableSetTypes {
		setMap[t] = make([]models.SetResponse, 0)
	}

	setMap["random"] = make([]models.SetResponse, 0)
	setMap["random"] = append(setMap["random"], models.SetResponse{
		Code: "RNG",
		Name: "Random Set",
	})
	sort.Strings(playableSetTypes)

	return setMap
}

func RefreshSetsInDB(context *gin.Context) {
	err := UpdateSets()
	if err != nil {
		context.JSON(500, gin.H{"error": "unexpected error"})
	}
}

func UpdateSets() error {
	log.Println("Refreshing sets")
	resp, err := http.Get("http://mtgjson.com/api/v5/AllPrintings.json")
	if err != nil {
		log.Println("Could not fetch the MTGJson allPrintings")
		return err
	}
	defer resp.Body.Close()

	allPrintings := new(mtgjson.AllPrintings)
	if err := json.NewDecoder(resp.Body).Decode(allPrintings); err != nil {
		log.Println("error while unmarshalling allPrintings", err)
		return err
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
			return err
		}
	}
	return nil
}

func RefreshSet(context *gin.Context) {
	log.Println("RefreshSetsInDB data")
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