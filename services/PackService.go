package services

import (
	"errors"
	"fmt"
	"math/rand"
	"mtgpoolservice/db"
	"mtgpoolservice/logging"
	"mtgpoolservice/models"
	"mtgpoolservice/models/entities"
	"time"
)

func MakePacks(sets []string) (packs []models.CardPool, err error) {
	for i := 0; i < len(sets); i++ {
		setCode := sets[i]

		set, err := db.GetSet(setCode)
		if err != nil {
			fmt.Println(err)
			logging.Warn(err)
			return nil, errors.New("set " + setCode + "does not exist")
		}

		pack, err := MakeRegularPack(set)
		if err != nil {
			fmt.Println(err)
			logging.Warn(err)
			return nil, errors.New("could not produce pack for " + setCode)
		}

		packs = append(packs, pack)
	}
	return
}

func CheckCubeList(req models.CubeListRequest) []string {
	_, missingCardNames := db.GetCardsByName(req.Cubelist[:])
	return missingCardNames
}

func MakeCubePacks(req *models.CubeDraftRequest) (packs []models.CardPool, err error) {
	cubeCards, missingCards := db.GetCardsByName(req.Cubelist)
	if len(missingCards) > 0 {
		return nil, fmt.Errorf("unknown cards", missingCards)
	}

	cubeCards.Shuffle()
	for i := 0; i < int(req.Players)*int(req.Packs); i++ {
		sliceLowerBound := i * int(req.PlayerPackSize)
		sliceUpperBound := sliceLowerBound + int(req.PlayerPackSize)
		slicedList := cubeCards[sliceLowerBound:sliceUpperBound]

		packs = append(packs, slicedList)
	}
	return
}

func MakeRegularPack(s *entities.Set) (models.CardPool, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	configuration, err := s.GetRandomConfiguration()
	if err != nil {
		return nil, err
	}

	protoCards := make([]entities.ProtoCard, 0)
	for _, confContent := range configuration.Contents {
		sheet, err := s.GetSheet(confContent.SheetName)
		if err != nil {
			return nil, err
		}

		randomCards := sheet.GetRandomCards(confContent.CardsNumber)
		protoCards = append(protoCards, randomCards...)
	}

	cards, err := getCards(s, protoCards)
	if err != nil {
		return nil, err
	}

	return cards, nil
}

func getCards(s *entities.Set, protoCards []entities.ProtoCard) (cardPool models.CardPool, err error) {
	for i, card := range protoCards {
		c, err := s.GetCard(card.UUID)
		if err != nil {
			return nil, err
		}

		cardPool.Add(c, protoCards[i].Foil)
	}
	return
}
