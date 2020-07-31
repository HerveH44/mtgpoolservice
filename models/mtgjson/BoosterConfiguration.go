package mtgjson

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	wr "mtgpoolservice/weighted"
)

type Configurations []PackConfiguration
type Sheets []Sheet

func (b *BoosterRule) GetSheet(name string) (*Sheet, error) {
	for _, sheet := range b.Sheets {
		if sheet.Name == name {
			return &sheet, nil
		}
	}

	return nil, errors.New("Could not find sheet with name " + name)
}

type BoosterRules []BoosterRule

type BoosterRule struct {
	SetID               string         `gorm:"primary_key"`
	Name                string         `gorm:"primary_key"`
	Boosters            Configurations `gorm:"foreignkey:SetID,BoosterRuleName;association_foreignkey:SetID,Name;PRELOAD:true"`
	Sheets              Sheets         `gorm:"foreignkey:SetID,BoosterRuleName;association_foreignkey:SetID,Name;PRELOAD:true"`
	BoostersTotalWeight uint
}

func (b *BoosterRules) UnmarshalJSON(data []byte) error {
	boosterMap := make(map[string]BoosterRule)
	err := json.Unmarshal(data, &boosterMap)
	if err != nil {
		return err
	}
	for name, config := range boosterMap {
		config.Name = name
		*b = append(*b, config)
	}

	return nil
}

type Sheet struct {
	ID              string
	Name            string
	SetID           string
	BoosterRuleName string
	BalanceColors   bool       `json:"balanceColors"`
	Foil            bool       `json:"foil"`
	TotalWeight     uint       `json:"totalWeight"`
	Cards           SheetCards `json:"cards;PRELOAD:true"`
}

func (user *Sheet) BeforeCreate(scope *gorm.Scope) error {
	sheet := scope.Value.(*Sheet)
	scope.SetColumn("ID", sheet.SetID+"_"+sheet.BoosterRuleName+"_"+sheet.Name)
	return nil
}

type ProtoCard struct {
	UUID string
	Foil bool
}

func (s *Sheet) GetRandomCards(cardsNumber uint) (ret []ProtoCard) {

	for i := 0; i < int(cardsNumber); i++ {
		choices := make([]wr.Choice, 0)
		for _, conf := range s.Cards {
			choices = append(choices, wr.Choice{
				Item:   conf,
				Weight: conf.Weight,
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
	Weight  uint
}

type SheetCards []SheetCard

func (b *SheetCards) UnmarshalJSON(data []byte) error {
	sheetsMap := make(map[MTGCardUUID]uint)
	if err := json.Unmarshal(data, &sheetsMap); err != nil {
		return fmt.Errorf("sheetcards: %w", err)
	}

	for uuid, weight := range sheetsMap {
		*b = append(*b, SheetCard{
			UUID:   uuid,
			Weight: weight,
		})
	}

	return nil
}

type SheetName = string
type MTGCardUUID = string

func asSha256(o interface{}) string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

type PackConfiguration struct {
	ID              string
	SetID           string
	BoosterRuleName string
	Weight          uint     `json:"weight"`
	Contents        Contents `json:"contents" gorm:"foreignkey:ConfigurationID;PRELOAD:true"`
}

func (user *PackConfiguration) BeforeCreate(scope *gorm.Scope) error {
	sheet := scope.Value.(*PackConfiguration)
	scope.SetColumn("ID", sheet.SetID+"_"+sheet.BoosterRuleName+"_"+asSha256(sheet.Contents))
	return nil
}

type ConfigurationContent struct {
	gorm.Model

	SheetName       string
	CardsNumber     uint
	ConfigurationID string
}

type Contents []ConfigurationContent

type ContentsMap map[SheetName]uint

func (b *Contents) UnmarshalJSON(data []byte) error {
	sheetsMap := make(ContentsMap)
	if err := json.Unmarshal(data, &sheetsMap); err != nil {
		return fmt.Errorf("contents: %w", err)
	}

	for sheetName, weight := range sheetsMap {
		*b = append(*b, ConfigurationContent{
			SheetName:   sheetName,
			CardsNumber: weight,
		})
	}

	return nil
}

type SheetsMap map[SheetName]Sheet

func (b *Sheets) UnmarshalJSON(data []byte) error {
	sheetsMap := make(SheetsMap)
	err := json.Unmarshal(data, &sheetsMap)
	if err != nil {
		return fmt.Errorf("sheets: %w", err)
	}
	for name, sheet := range sheetsMap {
		sheet.Name = name
		*b = append(*b, sheet)
	}

	return nil
}
