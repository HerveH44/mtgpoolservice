package entities

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm/dialects/postgres"
	"mtgpoolservice/utils"
)

type Set struct {
	Code               string `gorm:"primary_key"`
	Name               string
	Type               string
	ReleaseDate        string
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
