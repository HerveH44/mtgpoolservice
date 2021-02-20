package db

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/patrickmn/go-cache"
)

// SET REPO //

type SetRepository interface {
	GetRandomSet() (*Set, error)
	FindSet(setCode string) (*Set, error)
	FindAllSets() ([]*Set, error)
	FindLatestSet() (*Set, error)
	GetChaosSets(modernOnly bool) ([]*Set, error)
	SaveSets(sets *[]Set) error
}

type setRepo struct {
	db    *gorm.DB
	cache *cache.Cache
}

func NewSetRepository(db *gorm.DB) SetRepository {
	setCache := cache.New(6*time.Hour, 1*time.Hour)
	return &setRepo{db, setCache}
}

func (s *setRepo) GetChaosSets(modernOnly bool) (chaosSets []*Set, err error) {
	sets, err := s.FindAllSets()
	if err != nil {
		return
	}

	for _, set := range sets {
		if !set.IsExpansionOrCore() {
			continue
		}
		if modernOnly && !set.IsModern() {
			continue
		}
		chaosSets = append(chaosSets, set)
	}
	return chaosSets, nil
}

func (s *setRepo) SaveSets(sets *[]Set) error {
	log.Info("starting saving sets")
	//for _, set := range *sets {
	//	s.db.Create(&set)
	//}
	s.db.CreateInBatches(sets, 1)
	log.Info("finished saving sets")
	return nil
}

func (s *setRepo) GetRandomSet() (*Set, error) {
	var set = new(Set)

	err := s.db.
		Order("random()").
		Set("gorm:auto_preload", true).
		First(set).
		Error
	return set, err
}

func (s *setRepo) FindSet(setCode string) (*Set, error) {
	if cachedSet, found := s.cache.Get(setCode); found {
		return cachedSet.(*Set), nil
	}
	var set = new(Set)

	err := s.db.
		Where(" code = ?", setCode).
		Set("gorm:auto_preload", true).
		First(set).
		Error
	if err == nil {
		s.cache.SetDefault(setCode, set)
	}
	return set, err
}

func (s *setRepo) FindAllSets() (sets []*Set, err error) {
	err = s.db.
		Order("release_date DESC").
		Find(&sets).Error

	return
}

func (s *setRepo) FindLatestSet() (*Set, error) {
	var set = new(Set)
	err := s.db.
		Where("type in (?)", []string{"core", "expansion"}).
		Where("release_date <= now()").
		Order("release_date DESC").
		First(set).
		Error
	return set, err
}

// CARD REPO //

type CardRepository interface {
	GetCardWithName(name string) (*Card, error)
	GetCardWithNameAndSetInfos(name, setCode, number string) (*Card, error)
	GetCardsWithRarity(setCode, rarity string, number int) ([]*Card, error)
	GetRandomCardsFromSetsWithRarity(setCodes []string, rarity string, number int) ([]*Card, error)
}

type cardRepo struct {
	db    *gorm.DB
	cache *cache.Cache
}

var cardCache = cache.New(6*time.Hour, 1*time.Hour)
var cardWithSetInfosCache = cache.New(6*time.Hour, 1*time.Hour)
var unknownCardCache = cache.New(6*time.Hour, 1*time.Hour)

func NewCardRepository(db *gorm.DB) CardRepository {
	cardCache := cache.New(6*time.Hour, 1*time.Hour)
	return &cardRepo{db, cardCache}
}

func (c *cardRepo) GetCardWithName(name string) (*Card, error) {
	if foundCard, found := cardCache.Get(name); found {
		return foundCard.(*Card), nil
	}

	if foundUnknownCardError, found := unknownCardCache.Get(name); found {
		return nil, foundUnknownCardError.(error)
	}

	var err error
	var card = new(Card)
	var query *gorm.DB

	if isMultiCard := strings.ContainsAny(name, "/"); isMultiCard {
		query = c.db.Where("cubable = true AND name ILIKE ?", name)
	} else {
		query = c.db.Where("cubable = true AND face_name = ?", name)
	}

	if err = query.First(card).Error; err != nil {
		unknownCardCache.SetDefault(name, err)
		log.Debug("could not find card with name", name)
	} else {
		cardCache.SetDefault(name, card)
	}

	return card, err
}

func (c *cardRepo) GetCardWithNameAndSetInfos(name, setCode, number string) (*Card, error) {
	cacheCardRepresentation := fmt.Sprintf("%s|%s|%s", name, setCode, number)
	if foundCard, found := cardWithSetInfosCache.Get(cacheCardRepresentation); found {
		return foundCard.(*Card), nil
	}

	var err error
	var card = new(Card)
	var query *gorm.DB
	if isMultiCard := strings.ContainsAny(name, "/"); isMultiCard {
		query = c.db.
			Where("name ILIKE ?", name)
	} else {
		query = c.db.
			Where("face_name = ?", name)
	}

	err = query.
		Where("set_id = ?", strings.ToUpper(setCode)).
		Where("number ILIKE ?", number).
		First(card).
		Error

	if err != nil {
		log.Debug("could not find card with name", name, "setCode", setCode, "and number", number)
	}
	cardCache.SetDefault(cacheCardRepresentation, card)
	return card, err
}

func (c *cardRepo) GetCardsWithRarity(setCode, rarity string, number int) (cards []*Card, err error) {
	err = c.db.
		Raw("SELECT * from cards where set_id = ? AND rarity = ? ORDER BY random() LIMIT ?", setCode, rarity, number).
		Scan(&cards).
		Error
	return
}

func (c *cardRepo) GetRandomCardsFromSetsWithRarity(setCodes []string, rarity string, number int) (cards []*Card, err error) {
	err = c.db.
		Raw("SELECT * from cards where set_id in (?) AND rarity = ? ORDER BY random() LIMIT ?", setCodes, rarity, number).
		Scan(&cards).
		Error
	return
}

// VERSION REPO //

type VersionRepository interface {
	GetVersion() (Version, error)
	SaveVersion(version Version) error
}

type versionRepository struct {
	db *gorm.DB
}

func NewVersionRepository(db *gorm.DB) VersionRepository {
	return &versionRepository{db}
}

func (vr *versionRepository) GetVersion() (v Version, err error) {
	err = vr.db.Order("date DESC").First(&v).Error
	return
}

func (vr *versionRepository) SaveVersion(version Version) (err error) {
	return vr.db.Create(&version).Error
}
