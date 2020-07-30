package services

import (
	"errors"
	"fmt"
	"math/rand"
	"mtgpoolservice/common"
	"mtgpoolservice/models"
	"time"
)

func MakePacks(sets []string) (packs []models.Pack, err error) {
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

func MakePack(s *models.Set) (models.Pack, error) {
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

	protoCards := make([]models.ProtoCard, 0)
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
