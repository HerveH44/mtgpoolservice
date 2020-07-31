package services

import (
	"errors"
	"fmt"
	"math/rand"
	"mtgpoolservice/common"
	"mtgpoolservice/models"
	"mtgpoolservice/models/mtgjson"
	"time"
)

func MakePacks(sets []string) (packs []models.Pool, err error) {
	for i := 0; i < len(sets); i++ {
		setCode := sets[i]

		set, err := common.GetSet(setCode)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("set " + setCode + "does not exist")
		}

		pack, err := MakePack(&set)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("could not produce pack for " + setCode)
		}

		packs = append(packs, pack)
	}
	return
}

// Shuffle shuffles the array parameter in place
func Shuffle(a []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
}

func MakeCubePacks(req models.CubeSealedRequest) (packs []models.Pool, err error) {
	if int(req.PlayerPoolSize)*int(req.Players) > len(req.Cubelist) {
		return nil, fmt.Errorf("makecubepack: cube list too small")
	}

	Shuffle(req.Cubelist)

	for i := 0; i < int(req.Players); i++ {
		sliceLowerBound := i * int(req.PlayerPoolSize)
		sliceUpperBound := sliceLowerBound + int(req.PlayerPoolSize)
		slicedList := req.Cubelist[sliceLowerBound:sliceUpperBound]

		pack, err := common.GetCardsByName(slicedList)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("makecubepacks: could not produce pack")
		}

		packs = append(packs, pack)
	}
	return
}

func MakePack(s *mtgjson.Set) (models.Pool, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	boosterRule, err := s.GetDefaultBoosterRule()
	if err != nil {
		//TODO: make a default booster
		return nil, err
	}

	configuration, err := boosterRule.GetRandomConfiguration()
	if err != nil {
		return nil, err
	}

	protoCards := make([]mtgjson.ProtoCard, 0)
	for _, confContent := range configuration.Contents {
		sheet, err := boosterRule.GetSheet(confContent.SheetName)
		if err != nil {
			return nil, err
		}

		randomCards := sheet.GetRandomCards(confContent.CardsNumber)
		protoCards = append(protoCards, randomCards...)
	}

	cards, err := common.GetCards(protoCards)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found %s cards", cards)

	return cards, nil
}
