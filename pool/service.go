package pool

import (
	"errors"
	"fmt"
	"math/rand"
	"mtgpoolservice/db"
	"regexp"
	"runtime"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"

	"github.com/Jeffail/tunny"
)

type Service interface {
	checkCubeList(list []string) []string
	regularDraft(sets []string, playersNumber int) (packs []Pack, err error)
	regularSealed(sets []string, playersNumber int) (packs []Pack, err error)
	cubePacks(list []string, packSize, packsNumber int) (packs []Pack, err error)
	chaosDraft(modernOnly, totalChaos bool, packsNumber int) (packs []Pack, err error)
	chaosSealed(modernOnly, totalChaos bool, packsNumber, playersNumber int) (packs []Pack, err error)
}

type RandomPackResult struct {
	Pool  Pack
	Error error
}

type service struct {
	setRepo  db.SetRepository
	cardRepo db.CardRepository
}

func NewPackService(setRepo db.SetRepository, cardRepo db.CardRepository) Service {
	return &service{setRepo, cardRepo}
}

func (s *service) checkCubeList(list []string) []string {
	_, missingCardNames := s.checkList(list[:])
	return missingCardNames
}

func (s *service) regularDraft(sets []string, playersNumber int) (packs []Pack, err error) {
	for _, setCode := range sets {
		set, err := s.getSet(setCode)
		if err != nil {
			return nil, err
		}

		for i := 0; i < playersNumber; i++ {
			pack, err := s.makeRegularPack(set)
			if err != nil {
				log.Warn(err)
				return nil, errors.New("could not produce pack for " + set.Code)
			}

			packs = append(packs, pack)
		}
	}
	return
}

func (s *service) regularSealed(setCodes []string, players int) (pools []Pack, err error) {
	var sets = make([]*db.Set, 0)

	for _, setCode := range setCodes {
		set, err := s.getSet(setCode)
		if err != nil {
			return nil, err
		}
		sets = append(sets, set)
	}

	for i := 0; i < players; i++ {
		pool := Pack{}
		for _, set := range sets {
			pack, err := s.makeRegularPack(set)
			if err != nil {
				log.Warn(err)
				return nil, errors.New("could not produce pack for " + set.Code)
			}
			pool.AddPackContent(pack)
		}
		pools = append(pools, pool)
	}
	return
}

func (s *service) getSet(setCode string) (*db.Set, error) {
	var (
		set *db.Set
		err error
	)
	if setCode == "RNG" {
		set, err = s.setRepo.GetRandomSet()
		if err != nil {
			log.Warn(err)
			return nil, errors.New("server error: could not produce random set")
		}
	} else {
		set, err = s.setRepo.FindSet(setCode)
		if err != nil {
			log.Warn(err)
			return nil, errors.New("set " + setCode + "does not exist")
		}
	}
	return set, err
}

// TODO: Could be done by set directly
func (s *service) makeRegularPack(set *db.Set) (Pack, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	configuration, err := set.GetRandomConfiguration()
	if err != nil {
		return s.makeDefaultPack(set)
	}

	protoCards := make([]db.ProtoCard, 0)
	for _, confContent := range configuration.Contents {
		sheet, err := set.GetSheet(confContent.SheetName)
		if err != nil {
			return nil, err
		}

		randomCards := sheet.GetRandomCards(confContent.CardsNumber)
		protoCards = append(protoCards, randomCards...)
	}

	return getCards(set, protoCards)
}

// TODO: Could be done by set directly
func (s *service) makeDefaultPack(set *db.Set) (cards Pack, err error) {
	rares, err := s.cardRepo.GetCardsWithRarity(set.Code, "Rare", 1)
	if err != nil {
		return nil, err
	}
	unco, err := s.cardRepo.GetCardsWithRarity(set.Code, "Uncommon", 3)
	if err != nil {
		return nil, err
	}
	commons, err := s.cardRepo.GetCardsWithRarity(set.Code, "Common", 10)
	if err != nil {
		return nil, err
	}

	cards.AddCards(rares).
		AddCards(unco).
		AddCards(commons)
	return
}

func getCards(s *db.Set, protoCards []db.ProtoCard) (cardPool Pack, err error) {
	for i, card := range protoCards {
		c, err := s.GetCard(card.UUID)
		if err != nil {
			return nil, err
		}

		cardPool.Add(c, protoCards[i].Foil)
	}
	return
}

func (s *service) makeTotalChaosPack(sets []*db.Set) RandomPackResult {
	pack := Pack{}

	//TODO: get Mythic or Rare
	rares, err := s.getRandomCards(sets, 1, "Rare")
	if err != nil {
		return RandomPackResult{Error: err}
	}

	unco, err := s.getRandomCards(sets, 3, "Uncommon")
	if err != nil {
		return RandomPackResult{Error: err}
	}

	common, err := s.getRandomCards(sets, 10, "Common")
	if err != nil {
		return RandomPackResult{Error: err}
	}

	pack.AddCards(rares)
	pack.AddCards(unco)
	pack.AddCards(common)
	return RandomPackResult{Pool: pack}
}

func (s *service) getRandomCards(sets []*db.Set, number int, rarity string) ([]*db.Card, error) {
	setCodes := make([]string, len(sets))
	for _, set := range sets {
		setCodes = append(setCodes, set.Code)
	}
	return s.cardRepo.GetRandomCardsFromSetsWithRarity(setCodes, rarity, number)
}

func (s *service) makeRandomPack(sets []*db.Set) (ret RandomPackResult) {
	randomIndex := rand.Intn(len(sets))
	randomSet := sets[randomIndex]

	fullSet, err := s.setRepo.FindSet(randomSet.Code)
	if err != nil {
		return RandomPackResult{Error: err}
	}

	pack, err := s.makeRegularPack(fullSet)
	if err != nil {
		return RandomPackResult{Error: err}
	}

	return RandomPackResult{Pool: pack}
}

func (s *service) cubePacks(list []string, packSize, packsNum int) (packs []Pack, err error) {
	cubeCards, missingCards := s.checkList(list[:])
	if len(missingCards) > 0 {
		return nil, fmt.Errorf("unknown cards: %s", missingCards)
	}

	cubeCards.Shuffle()
	for i := 0; i < packsNum; i++ {
		sliceLowerBound := i * packSize
		sliceUpperBound := sliceLowerBound + packSize
		slicedList := cubeCards[sliceLowerBound:sliceUpperBound]

		packs = append(packs, slicedList)
	}
	return
}

func (s *service) chaosDraft(modernOnly, totalChaos bool, packsNumber int) (packs []Pack, err error) {
	sets, err := s.setRepo.GetChaosSets(modernOnly)
	if err != nil {
		return
	}

	numCPUs := runtime.NumCPU()
	pool := tunny.NewFunc(numCPUs, func(payload interface{}) interface{} {
		sets := payload.([]*db.Set)
		if !totalChaos {
			return s.makeRandomPack(sets)
		}
		return s.makeTotalChaosPack(sets)
	})
	defer pool.Close()

	for i := 0; i < packsNumber; i++ {
		result := pool.Process(sets).(RandomPackResult)
		if result.Error != nil {
			i--
			continue
		}
		packs = append(packs, result.Pool)
	}

	return
}

func (s *service) chaosSealed(modernOnly, totalChaos bool, packsNumber, playersNumber int) (packs []Pack, err error) {
	sets, err := s.setRepo.GetChaosSets(modernOnly)
	if err != nil {
		return
	}

	for i := 0; i < playersNumber; i++ {
		pool := Pack{}
		for j := 0; j < packsNumber; j++ {
			var result RandomPackResult
			if !totalChaos {
				result = s.makeRandomPack(sets)
			} else {
				result = s.makeTotalChaosPack(sets)
			}
			pool.AddPackContent(result.Pool)
		}
		packs = append(packs, pool)
	}

	return
}

func (s *service) checkList(names []string) (cr Pack, missingCardNames []string) {
	jobs := make(chan string, len(names))
	missingCards := make(chan string)
	cards := make(chan db.Card)
	var wg sync.WaitGroup

	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go s.getCard(jobs, missingCards, cards, &wg)
	}

	go addToMissingCards(missingCards, &missingCardNames)
	go addToCardPool(cards, &cr)

	for _, name := range names {
		jobs <- name
	}

	close(jobs)
	wg.Wait()
	close(cards)
	close(missingCards)
	return
}

func (s *service) getCard(jobs <-chan string, missingCards chan<- string, foundCards chan<- db.Card, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		var card = regexp.MustCompile(`^(.*?)(?: +\((\w+)(?: +(\w+))?\))? *$`)
		parts := card.FindStringSubmatch(j)
		var c *db.Card
		var err error

		if len(parts) == 4 {
			c, err = s.cardRepo.GetCardWithNameAndSetInfos(parts[1], parts[2], parts[3])
		} else {
			c, err = s.cardRepo.GetCardWithName(j)
		}
		if err != nil {
			missingCards <- j
		} else {
			foundCards <- *c
		}
	}
}

func addToCardPool(cards chan db.Card, cr *Pack) {
	for c := range cards {
		cr.Add(&c, false)
	}
}

func addToMissingCards(missingCards <-chan string, missingCardNames *[]string) {
	for c := range missingCards {
		*missingCardNames = append(*missingCardNames, c)
	}
}

type cardDTO struct {
	db.Card

	Id   string `json:"cardId"`
	Foil bool   `json:"foil"`
}

type Pack []cardDTO

// Shuffle shuffles the array parameter in place
func (c *Pack) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(*c), func(i, j int) { (*c)[i], (*c)[j] = (*c)[j], (*c)[i] })
}

func (c *Pack) Add(card *db.Card, isFoil bool) {
	cardResponse := cardDTO{
		Card: *card,
		Id:   uuid.New().String(),
		Foil: isFoil,
	}
	*c = append(*c, cardResponse)
}

func (c *Pack) AddCards(cards []*db.Card) *Pack {
	for _, card := range cards {
		c.Add(card, false)
	}
	return c
}

func (c *Pack) AddPackContent(pack Pack) {
	for _, card := range pack {
		*c = append(*c, card)
	}
}
