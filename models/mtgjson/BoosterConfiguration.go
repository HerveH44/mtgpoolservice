package mtgjson

import (
	"encoding/json"
	"fmt"
)

type BoosterRule struct {
	Sheets              map[string]Sheet
	Boosters            []PackConfiguration
	BoostersTotalWeight int
}

type Sheets = []Sheet

//func (b *Sheets) UnmarshalJSON(data []byte) error {
//	sheetsMap := make(map[string]Sheet)
//	err := json.Unmarshal(data, &sheetsMap)
//	if err != nil {
//		return fmt.Errorf("sheets: %w", err)
//	}
//	for name, sheet := range sheetsMap {
//		sheet.Name = name
//		*b = append(*b, sheet)
//	}
//
//	return nil
//}

type Sheet struct {
	Name          string
	BalanceColors bool       `json:"balanceColors"`
	Foil          bool       `json:"foil"`
	TotalWeight   uint       `json:"totalWeight"`
	Cards         SheetCards `json:"cards"`
}

type SheetCard struct {
	UUID   string
	Weight int
}

type SheetCards []SheetCard

func (b *SheetCards) UnmarshalJSON(data []byte) error {
	sheetsMap := make(map[MTGCardUUID]int)
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

type MTGCardUUID = string

type PackConfiguration struct {
	Weight   int      `json:"weight"`
	Contents Contents `json:"contents"`
}

type Contents []ConfigurationContent

type ConfigurationContent struct {
	SheetName   string
	CardsNumber int
}

func (b *Contents) UnmarshalJSON(data []byte) error {
	sheetsMap := make(map[string]int)
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
