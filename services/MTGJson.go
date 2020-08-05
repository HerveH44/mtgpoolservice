package services

import (
	"encoding/json"
	"fmt"
	"log"
	database "mtgpoolservice/db"
	"mtgpoolservice/models/mtgjson"
	"mtgpoolservice/utils"
	"net/http"
	"sort"
	"strings"
	"time"
)

var mtgJsonEndpoint = "http://mtgjson.com/api/v5/"
var notCubableSetTypes = []string{"box", "duel_deck", "masterpiece", "memorabilia", "promo", "spellbook"}

func UpdateSets() error {
	log.Println("Refreshing sets")
	resp, err := http.Get(mtgJsonEndpoint + "AllPrintings.json")
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

func UpdateSet(setCode string) error {
	resp, err := http.Get(fmt.Sprint(mtgJsonEndpoint, strings.ToUpper(setCode), ".json"))
	if err != nil {
		return fmt.Errorf("RefreshSet: Could not fetch the MTGJson set with code %s", setCode)
	}
	defer resp.Body.Close()

	monoSet := new(mtgjson.MonoSet)
	if err := json.NewDecoder(resp.Body).Decode(monoSet); err != nil {
		return fmt.Errorf("RefreshSet: Error while unmarshalling monoSet %w", err)
	}

	log.Println("main: saving set ", monoSet.Data.Name)
	entity := mtgjson.MapMTGJsonSetToEntity(monoSet.Data, notCubableFunc)
	if err := database.GetDB().Save(&entity).Error; err != nil {
		return fmt.Errorf("RefreshSet: could not save the set %s - %w", monoSet.Data.Name, err)
	}
	return nil
}

func notCubableFunc(string, *mtgjson.Card) bool {
	return false
}
