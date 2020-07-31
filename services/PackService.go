package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"mtgpoolservice/db"
	"mtgpoolservice/logging"
	"mtgpoolservice/models"
	"mtgpoolservice/models/entities"
	"time"
)

func MakePacks(sets []string) (packs []models.Pool, err error) {
	for i := 0; i < len(sets); i++ {
		setCode := sets[i]

		set, err := db.GetSet(setCode)
		if err != nil {
			fmt.Println(err)
			logging.Warn(err)
			return nil, errors.New("set " + setCode + "does not exist")
		}

		pack, err := MakePack(set)
		if err != nil {
			fmt.Println(err)
			logging.Warn(err)
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

		pack, err := db.GetCardsByName(slicedList)
		if err != nil {
			fmt.Println(err)
			logging.Warn(err)
			return nil, fmt.Errorf("makecubepacks: could not produce pack")
		}

		packs = append(packs, pack)
	}
	return
}

func MakePack(s *entities.Set) (models.Pool, error) {
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

func getCards(s *entities.Set, protoCards []entities.ProtoCard) (cr []models.CardResponse, err error) {
	for i, card := range protoCards {
		c, err := s.GetCard(card.UUID)
		if err != nil {
			return nil, err
		}
		cr = append(cr, models.CardResponse{
			Card: c,
			Id:   uuid.New().String(),
			Foil: protoCards[i].Foil,
		})
	}
	return
}
