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
	"sort"
	"strings"
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
func Shuffle(a []models.CardResponse) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
}

func CheckCubeList(req models.CubeListRequest) (invalidCards []string) {
	invalidCards = db.CheckCubeCards(req.Cubelist[:])
	return
}

func MakeCubePacks(req models.CubeDraftRequest) (packs []models.Pool, err error) {
	if int(req.PlayerPackSize)*int(req.Players)*int(req.Packs) > len(req.Cubelist) {
		return nil, fmt.Errorf("MakeCubePacks: cube list too small")
	}

	cubeCards, err := db.GetCardsByName(req.Cubelist)
	if len(cubeCards) != len(req.Cubelist) {
		missingCards := GetMissingCards(req.Cubelist[:], cubeCards[:])
		fmt.Println("MakeCubePacks: could not find all the cards. Expected ", len(req.Cubelist), " cards but found ", len(cubeCards), "cards instead")
		fmt.Println(missingCards)
	}

	Shuffle(cubeCards[:])
	if err != nil {
		fmt.Println(err)
		logging.Warn(err)
		return nil, fmt.Errorf("MakeCubePacks: could not produce pack")
	}
	for i := 0; i < int(req.Players)*int(req.Packs); i++ {
		sliceLowerBound := i * int(req.PlayerPackSize)
		sliceUpperBound := sliceLowerBound + int(req.PlayerPackSize)
		slicedList := cubeCards[sliceLowerBound:sliceUpperBound]

		packs = append(packs, slicedList)
	}
	return
}

func GetMissingCards(cubeList []string, fetchedCards []models.CardResponse) (ret []string) {
	sort.Strings(cubeList)

	for _, card := range fetchedCards {
		index := sort.SearchStrings(cubeList, strings.ToLower(card.Name))
		if index >= len(cubeList) || cubeList[index] != strings.ToLower(card.Name) {
			ret = append(ret, card.Name)
		}
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
