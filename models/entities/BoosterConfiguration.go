package entities

import (
	"mtgpoolservice/utils"
)

type Sheet struct {
	ID            string `gorm:"primary_key"`
	SetID         string
	Name          string
	BalanceColors bool
	Foil          bool
	TotalWeight   uint
	Cards         []SheetCard
}

func (s *Sheet) GetRandomCards(cardsNumber int) (ret []ProtoCard) {
	choices := make([]utils.Choice, 0)
	for _, conf := range s.Cards {
		choices = append(choices, utils.NewChoice(conf, uint(conf.Weight)))
	}

	chooser := utils.NewChooser(choices...)
	for i := 0; i < cardsNumber; i++ {
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
