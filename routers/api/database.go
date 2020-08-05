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
	"mtgpoolservice/utils"
	"net/http"
	"sort"
	"strings"
	"time"
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
	context.JSON(200, gin.H{"success": "ok"})
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
	ordereredSetsCode := orderSetsByDate(&allPrintings.Data)
	isCubable := buildIsCubableFunc(ordereredSetsCode)
	for setName, set := range allPrintings.Data {
		i++
		fmt.Printf("%d/%d - saving set %s\n", i, setsNumber, setName)
		entity := mtgjson.MapMTGJsonSetToEntity(set, isCubable)
		if err := database.GetDB().Save(&entity).Error; err != nil {
			fmt.Printf("could not save the card %s - %s\n", setName, err)
			return err
		}
	}
	return nil
}

func buildIsCubableFunc(orderedSetCodes map[string]int) func(string, *mtgjson.Card) bool {
	return func(setCode string, card *mtgjson.Card) bool {
		if len(card.Variations) > 0 && (card.BorderColor == "borderless" || card.IsStarter || utils.Include(card.FrameEffects, "extendedart")) {
			return false
		}
		printings := card.Printings
		if len(printings) < 2 {
			return true
		}

		sort.SliceStable(printings, func(i, j int) bool {
			setCode1 := printings[i]
			setCode2 := printings[j]
			set1Index, ok := orderedSetCodes[setCode1]
			if !ok {
				return false
			}
			set2Index, ok := orderedSetCodes[setCode2]
			if !ok {
				return true
			}
			return set1Index < set2Index
		})

		if setCode == printings[0] {
			return true
		}

		return false
	}
}

var notCubableSetTypes = []string{"box", "duel_deck", "masterpiece", "memorabilia", "promo", "spellbook"}

func orderSetsByDate(m *map[string]mtgjson.MTGJsonSet) map[string]int {
	r := make([]string, 0)
	for setCode, set := range *m {
		if !utils.Include(notCubableSetTypes, set.Type) {
			r = append(r, setCode)
		}
	}
	myMap := *m
	sort.SliceStable(r, func(i, j int) bool {
		releaseDate1 := myMap[r[i]].ReleaseDate
		releaseDate2 := myMap[r[j]].ReleaseDate
		return time.Time(releaseDate1).After(time.Time(releaseDate2))
	})

	ret := make(map[string]int)
	for i, code := range r {
		ret[code] = i
	}
	return ret
}

func RefreshSet(context *gin.Context) {
	setCode := context.Param("setCode")
	if setCode == "" {
		context.JSON(500, "unexpected error: Could not fetch the MTGJson monoSet")
	}
	log.Println("RefreshSet data")
	resp, err := http.Get(fmt.Sprint("http://mtgjson.com/api/v5/", strings.ToUpper(setCode), ".json"))
	if err != nil {
		log.Println("RefreshSet: Could not fetch the MTGJson set with code", setCode)
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
	entity := mtgjson.MapMTGJsonSetToEntity(monoSet.Data, func(s string, i *mtgjson.Card) bool {
		return false
	})
	if err := database.GetDB().Save(&entity).Error; err != nil {
		log.Fatal("main: could not save the set", monoSet.Data.Name, err)

	}
}
