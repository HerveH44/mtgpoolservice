package entities

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"mtgpoolservice/utils"
	wr "mtgpoolservice/weighted"
)

type BoosterRule struct {
	ID                  string              `gorm:primary_key`
	Sheets              []Sheet             `gorm:"foreignkey:SetID;PRELOAD:true"`
	PackConfigurations  []PackConfiguration `gorm:"foreignkey:SetID;PRELOAD:true"`
	BoostersTotalWeight int
}

func (b *BoosterRule) GetSheet(name string) (*Sheet, error) {
	for _, sheet := range b.Sheets {
		if sheet.Name == name {
			return &sheet, nil
		}
	}

	return nil, errors.New("Could not find sheet with name " + name)
}

func (r *BoosterRule) GetRandomConfiguration() (*PackConfiguration, error) {
	configurations := r.PackConfigurations
	if len(configurations) == 0 {
		return nil, fmt.Errorf("BoosterRule.GetRandomConfiguration: Did not find any booster rule for %s", r.ID)
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

type Sheet struct {
	ID            string `gorm:"primary_key"`
	SetID         string
	Name          string
	BalanceColors bool
	Foil          bool
	TotalWeight   uint
	Cards         []SheetCard
}

// For internal use
type ProtoCard struct {
	UUID string
	Foil bool
}

func (s *Sheet) GetRandomCards(cardsNumber int) (ret []ProtoCard) {

	for i := 0; i < cardsNumber; i++ {
		choices := make([]wr.Choice, 0)
		for _, conf := range s.Cards {
			choices = append(choices, wr.Choice{
				Item:   conf,
				Weight: uint(conf.Weight),
			})
		}

		chooser := wr.NewChooser(choices...)
		pick := chooser.Pick().(SheetCard)

		ret = append(ret, ProtoCard{
			UUID: pick.UUID,
			Foil: s.Foil,
		})
	}

	return
}

type SheetCard struct {
	SheetID string `gorm:"primary_key"`
	UUID    string `gorm:"primary_key"`
	Weight  int
}

type PackConfiguration struct {
	ID              string
	SetID           string
	BoosterRuleName string
	Weight          int
	Contents        Contents `gorm:"foreignkey:ConfigurationID;PRELOAD:true"`
}

func (user *PackConfiguration) BeforeCreate(scope *gorm.Scope) error {
	sheet := scope.Value.(*PackConfiguration)
	scope.SetColumn("ID", sheet.SetID+"_"+sheet.BoosterRuleName+"_"+utils.AsSha256(sheet.Contents))
	return nil
}

type Contents []ConfigurationContent

type ConfigurationContent struct {
	ConfigurationID string
	SheetName       string
	CardsNumber     int
}
