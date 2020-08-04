package db

import (
	"fmt"
	"log"
	"mtgpoolservice/models"
	"mtgpoolservice/setting"
	"sort"
	"testing"
)

func init() {
	setting.Setup()
}

func TestGetSet(t *testing.T) {
	Init()
	set, err := GetSet("ISD")

	if err != nil {
		t.Error(err)
	}

	if set.Code != "ISD" {
		t.Error("expected ISD set")
	}
}

func TestGetSets(t *testing.T) {
	Init()
	sets, err := GetSets()

	if err != nil {
		t.Error(err)
	}

	setMap := make(map[string][]models.SetResponse)
	playableSetTypes := []string{"core", "expansion", "draft_innovation", "funny", "starter", "masters"}
	for _, t := range playableSetTypes {
		setMap[t] = make([]models.SetResponse, 0)
	}
	sort.Strings(playableSetTypes)
	for _, set := range sets {
		setType := set.Type
		i := sort.SearchStrings(playableSetTypes, setType)
		if i >= len(playableSetTypes) || playableSetTypes[i] != setType {
			fmt.Println("Found a bad Set!", set.Name, set.Type)
			continue
		}
		setMap[setType] = append(setMap[setType], models.SetResponse{
			Code: set.Code,
			Name: set.Name,
		})
	}
	setMap["random"] = make([]models.SetResponse, 0)
	setMap["random"] = append(setMap["random"], models.SetResponse{
		Code: "RNG",
		Name: "Random Set",
	})
	log.Println(setMap)
}
