package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"mtgpoolservice/utils"
	"time"

	"golang.org/x/mod/semver"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type Version struct {
	Date            time.Time `gorm:"Type:date"`
	SemanticVersion string
}

func (v *Version) IsNewer(lastVersion *Version) bool {
	if v.Date.After(lastVersion.Date) {
		return true
	}
	compareVersion := semver.Compare(v.SemanticVersion, lastVersion.SemanticVersion)
	if compareVersion == 1 {
		return true
	}
	return false
}

type Set struct {
	Code               string `gorm:"primary_key"`
	Name               string
	Type               string
	ReleaseDate        time.Time `gorm:"Type:date"`
	BaseSetSize        int
	Cards              []Card  `gorm:"foreignkey:SetID;PRELOAD:true"`
	Sheets             []Sheet `gorm:"foreignkey:SetID;PRELOAD:true"`
	PackConfigurations postgres.Jsonb

	parsedConfigurations []PackConfiguration `gorm:"-"`
	cardsByUuid          map[string]*Card    `gorm:"-"`
}

func (s *Set) GetSheet(name string) (*Sheet, error) {
	for _, sheet := range s.Sheets {
		if sheet.Name == name {
			return &sheet, nil
		}
	}

	return nil, errors.New("Could not find sheet with name " + name)
}

func (s *Set) GetRandomConfiguration() (*PackConfiguration, error) {
	configurations := s.getPackConfigurations()
	if len(configurations) == 0 {
		return nil, fmt.Errorf("BoosterRule.GetRandomConfiguration: Did not find any booster rule for %s", s.Code)
	}

	choices := make([]utils.Choice, 0)
	for _, conf := range configurations {
		choices = append(choices, utils.NewChoice(conf, uint(conf.Weight)))
	}

	chooser := utils.NewChooser(choices...)
	pick := chooser.Pick().(PackConfiguration)

	return &pick, nil
}

func (s *Set) getPackConfigurations() (ret []PackConfiguration) {
	if len(s.parsedConfigurations) > 0 {
		return s.parsedConfigurations
	}
	if err := json.Unmarshal(s.PackConfigurations.RawMessage, &ret); err != nil {
		fmt.Println("Set.getPackConfigurations()", "Could not parse json content for pack confs ", s.Code)
	}
	s.parsedConfigurations = ret
	return
}

func (s *Set) GetCard(uuid string) (c *Card, err error) {
	// Load cardsByUuid map
	if len(s.cardsByUuid) == 0 {
		s.cardsByUuid = make(map[string]*Card)
		for i, card := range s.Cards {
			s.cardsByUuid[card.UUID] = &s.Cards[i]
		}
	}

	c = s.cardsByUuid[uuid]
	if c == nil {
		return nil, fmt.Errorf("set[%s] card [%s] was not found in DB", s.Code, uuid)
	}
	return
}

func (s *Set) IsExpansionOrCore() bool {
	return s.Type == "expansion" || s.Type == "core"
}

var modernTime = time.Date(2003, 07, 25, 0, 0, 0, 0, time.UTC)

func (s *Set) IsModern() bool {
	return modernTime.Before(s.ReleaseDate)
}

type Card struct {
	UUID              string     `json:"uuid" gorm:"primary_key"`
	Name              string     `json:"name"`
	FaceName          string     `json:"-" gorm:"index:cubable_idx"`
	Color             string     `json:"color"`
	SetID             string     `json:"setCode" gorm:"index:pack_idx"`
	ConvertedManaCost int        `json:"cmc"`
	Number            string     `json:"number"`
	Type              string     `json:"type"`
	ManaCost          string     `json:"manaCost"`
	Rarity            string     `json:"rarity" gorm:"index:pack_idx"`
	URL               string     `json:"url"`
	ScryfallID        Identifier `json:"identifiers"` // Must be rename as identifiers with inner scryfallId
	Layout            string     `json:"layout"`
	IsDoubleFaced     bool       `json:"isDoubleFaced"`
	FlippedCardURL    string     `json:"flippedCardURL"`
	FlippedIsBack     bool       `json:"flippedIsBack"`
	FlippedNumber     string     `json:"flippedNumber"`
	Text              string     `json:"text"`
	Loyalty           string     `json:"loyalty"`
	Power             string     `json:"power"`
	Toughness         string     `json:"toughness"`
	Side              string     `json:"side"`          // no use ?
	IsAlternative     bool       `json:"isAlternative"` // no use ?
	Cubable           bool       `json:"-" gorm:"index:cubable_idx"`
}

type Identifier string

func (i *Identifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Inner string `json:"scryfallId"`
	}{
		Inner: string(*i),
	})
}

type Sheet struct {
	ID            string `gorm:"primary_key"`
	SetID         string
	Name          string
	BalanceColors bool
	Foil          bool
	TotalWeight   uint
	SheetCards    []SheetCard
}

func (s *Sheet) GetRandomCards(cardsNumber int) (ret []ProtoCard) {
	for i := 0; i < cardsNumber; i++ {

		choices := make([]utils.Choice, 0)
	OUTER:
		for _, conf := range s.SheetCards {
			for _, c := range ret {
				if c.UUID == conf.UUID {
					continue OUTER
				}
			}
			choices = append(choices, utils.NewChoice(conf, uint(conf.Weight)))
		}

		chooser := utils.NewChooser(choices...)
		pick := chooser.Pick().(SheetCard)
		ret = append(ret, ProtoCard{
			UUID: pick.UUID,
			Foil: s.Foil,
		})
	}

	return
}

// For internal use
type ProtoCard struct {
	UUID string
	Foil bool
}

type SheetCard struct {
	SheetID string `gorm:"primary_key"`
	UUID    string `gorm:"primary_key"`
	Weight  int
}

type PackConfiguration struct {
	Weight   int
	Contents Contents
}

type Contents []ConfigurationContent

type ConfigurationContent struct {
	SheetName   string
	CardsNumber int
}
