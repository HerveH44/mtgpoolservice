package services

import (
	"errors"
	"fmt"
	"github.com/Jeffail/tunny"
	"math/rand"
	"mtgpoolservice/db"
	"mtgpoolservice/logging"
	"mtgpoolservice/models"
	"mtgpoolservice/models/entities"
	"runtime"
	"time"
)

func MakeRegularPacks(sets []string) (packs []*models.CardPool, err error) {
	for i := 0; i < len(sets); i++ {
		setCode := sets[i]

		set, err := db.GetSet(setCode)
		if err != nil {
			fmt.Println(err)
			logging.Warn(err)
			return nil, errors.New("set " + setCode + "does not exist")
		}

		pack, err := makeRegularPack(set)
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

func makeRegularPack(s *entities.Set) (*models.CardPool, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	configuration, err := s.GetRandomConfiguration()
	if err != nil {
		return makeDefaultPack(s)
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

	return &cards, nil
}

func makeDefaultPack(s *entities.Set) (cards *models.CardPool, err error) {
	cards = &models.CardPool{}
	rares, err := db.GetCardsWithRarity(s.Code, 1, "Rare")
	if err != nil {
		return nil, err
	}
	unco, err := db.GetCardsWithRarity(s.Code, 3, "Uncommon")
	if err != nil {
		return nil, err
	}
	commons, err := db.GetCardsWithRarity(s.Code, 10, "Common")
	if err != nil {
		return nil, err
	}

	cards.AddCards(&rares)
	cards.AddCards(&unco)
	cards.AddCards(&commons)
	return
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

func MakeChaosPacks(req *models.ChaosRequest) (packs []models.CardPool, err error) {
	sets, err := getChaosSets(req.Modern)
	if err != nil {
		return
	}
	if !req.TotalChaos {
		numCPUs := runtime.NumCPU()

		pool := tunny.NewFunc(numCPUs, func(payload interface{}) interface{} {
			sets := payload.(*[]entities.Set)
			return makeRandomPack(sets)
		})
		defer pool.Close()
		for i := 0; i < int(req.Players*req.Packs); i++ {
			result := pool.Process(sets).(RandomPackResult)
			if result.Error != nil {
				i--
				continue
			}
			packs = append(packs, *result.Pool)
		}
	}

	return
}

type RandomPackResult struct {
	Pool  *models.CardPool
	Error error
}

func makeRandomPack(sets *[]entities.Set) (ret RandomPackResult) {
	randomIndex := rand.Intn(len(*sets))
	randomSet := (*sets)[randomIndex]

	fullSet, err := db.GetSet(randomSet.Code)
	if err != nil {
		return RandomPackResult{Error: err}
	}

	pack, err := makeRegularPack(fullSet)
	if err != nil {
		return RandomPackResult{Error: err}
	}

	return RandomPackResult{Pool: pack}
}
