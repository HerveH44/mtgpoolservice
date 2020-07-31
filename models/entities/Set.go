package entities

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm/dialects/postgres"
	wr "mtgpoolservice/weighted"
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

	choices := make([]wr.Choice, 0)
	for _, conf := range configurations {
		choices = append(choices, wr.Choice{
			Item:   conf,
			Weight: uint(conf.Weight),
		})
	}

	chooser := wr.NewChooser(choices...)
	pick := chooser.Pick().(PackConfiguration)

	return &pick, nil
}

func (s *Set) getPackConfigurations() (ret []PackConfiguration) {
	if err := json.Unmarshal(s.PackConfigurations.RawMessage, &ret); err != nil {
		fmt.Println("Set.getPackConfigurations()", "Could not parse json content for pack confs ", s.Code)
	}
	return
}
