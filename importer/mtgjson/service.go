package mtgjson

import (
	"encoding/json"
	"fmt"
	"log"
	"mtgpoolservice/utils"
	"net/http"
	"sort"
	"strings"
)

type MTGJsonService interface {
	DownloadVersion() (Version, error)
	DownloadSets(onSet func(set *MTGJsonSet, isCubable IsCubable)) error
	DownloadSet(setCode string) (*MTGJsonSet, error)
}

type IsCubable func(string, *Card) bool

var notCubableSetTypes = []string{"box", "duel_deck", "masterpiece", "memorabilia", "promo", "spellbook"}

type mtgJsonService struct {
	endpoint string
}

func NewMTGJsonService(endpoint string) MTGJsonService {
	return &mtgJsonService{endpoint: endpoint}
}

func (m *mtgJsonService) DownloadVersion() (version Version, err error) {
	log.Println("Check MTGJSON remote version")

	resp, err := http.Get(m.endpoint + "Meta.json")
	if err != nil {
		log.Println("Could not fetch the MTGJson version")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&version); err != nil {
		log.Println("error while unmarshalling version", err)
	}

	return
}

func (m *mtgJsonService) DownloadSet(setCode string) (*MTGJsonSet, error) {
	resp, err := http.Get(fmt.Sprint(m.endpoint, strings.ToUpper(setCode), ".json"))
	if err != nil {
		return nil, fmt.Errorf("RefreshSet: Could not fetch the MTGJson set with code %s", setCode)
	}
	defer resp.Body.Close()

	monoSet := new(MonoSet)
	if err := json.NewDecoder(resp.Body).Decode(monoSet); err != nil {
		return nil, fmt.Errorf("RefreshSet: Error while unmarshalling monoSet %w", err)
	}

	return &monoSet.Data, nil
}

func (m *mtgJsonService) DownloadSets(onSet func(set *MTGJsonSet, isCubable IsCubable)) error {
	log.Println("Refreshing sets")
	resp, err := http.Get(m.endpoint + "AllPrintings.json")
	if err != nil {
		log.Println("Could not fetch the MTGJson allPrintings")
		return err
	}
	defer resp.Body.Close()

	allPrintings := new(AllPrintings)
	if err := json.NewDecoder(resp.Body).Decode(allPrintings); err != nil {
		log.Println("error while unmarshalling allPrintings", err)
		return err
	}
	var setsNumber = len(allPrintings.Data)
	log.Println("sets found", setsNumber)

	ordereredSetsCode := orderSetsByDate(&allPrintings.Data)
	isCubable := buildIsCubableFunc(ordereredSetsCode)

	for _, set := range allPrintings.Data {
		onSet(&set, isCubable)
	}

	return nil
}

func buildIsCubableFunc(orderedSetCodes map[string]int) func(string, *Card) bool {
	return func(setCode string, card *Card) bool {
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

func orderSetsByDate(m *map[string]MTGJsonSet) map[string]int {
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
		return releaseDate1.Before(&releaseDate2)
	})

	ret := make(map[string]int)
	for i, code := range r {
		ret[code] = i
	}
	return ret
}
