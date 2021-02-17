package mtgjson

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Meta struct {
	Date    VersionDate `json:"date"`
	Version string      `json:"version"`
}

type VersionDate time.Time

func (v *VersionDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*v = VersionDate(t)
	return nil
}

func (v *VersionDate) Before(v2 *VersionDate) bool {
	return time.Time(*v).Before(time.Time(*v2))
}

type Version struct {
	Data struct {
		Date    VersionDate `json:"date"`
		Version string      `json:"version"`
	} `json:"data"`

	Meta struct {
		Date    VersionDate `json:"date"`
		Version string      `json:"version"`
	} `json:"meta"`
}

/**
Représentation de l'api MTGJson
comme on trouve dans /api/v5/AllPrintings.json
*/
type AllPrintings struct {
	Data map[string]MTGJsonSet `json:"data"`
	Meta Meta                  `json:"meta"`
}

/**
Représentation d'un Set MTGJson
comme on trouve dans /api/v5/ISD.json
*/
type MonoSet struct {
	Data MTGJsonSet `json:"data"`
	Meta Meta       `json:"meta"`
}

type MTGJsonSet struct {
	Code         string      `json:"code"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Block        string      `json:"block"`
	ReleaseDate  VersionDate `json:"releaseDate"`
	BaseSetSize  int         `json:"baseSetSize"`
	TotalSetSize int         `json:"totalSetSize"`
	Booster      struct {
		Default BoosterRule `json:"default"` // on ne veut que le booster de défaut
	} `json:"booster"`
	Cards []Card `json:"cards"`
}

type Card struct {
	UUID              string   `json:"uuid"`
	Name              string   `json:"name"`
	FaceName          string   `json:"faceName,omitempty"`
	ColorIdentity     []string `json:"colorIdentity"`
	Colors            []string `json:"colors"`
	ColorIndicator    []string `json:"colorIndicator,omitempty"`
	ManaCost          string   `json:"manaCost,omitempty"`
	ConvertedManaCost float32  `json:"convertedManaCost"`
	FrameVersion      string   `json:"frameVersion"`
	Loyalty           string   `json:"loyalty"`
	Identifiers       struct {
		ScryfallID             string `json:"scryfallId"`
		ScryfallIllustrationID string `json:"scryfallIllustrationId"`
		ScryfallOracleID       string `json:"scryfallOracleId"`
	} `json:"identifiers,omitempty"`
	Layout        string   `json:"layout"`
	Number        string   `json:"number"`
	Power         string   `json:"power,omitempty"`
	Toughness     string   `json:"toughness,omitempty"`
	Printings     []string `json:"printings"`
	Rarity        string   `json:"rarity"`
	Text          string   `json:"text,omitempty"`
	Type          string   `json:"type"`
	Types         []string `json:"types"`
	Subtypes      []string `json:"subtypes"`
	Supertypes    []string `json:"supertypes"`
	IsReprint     bool     `json:"isReprint,omitempty"`
	IsAlternative bool     `json:"isAlternative,omitempty"`
	FrameEffects  []string `json:"frameEffects,omitempty"`
	OtherFaceIds  []string `json:"otherFaceIds,omitempty"`
	Side          string   `json:"side,omitempty"`
	Variations    []string `json:"variations,omitempty"`
	IsPromo       bool     `json:"isPromo"`
	IsStarter     bool     `json:"isStarter"`
	BorderColor   string   `json:"borderColor"`
}

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
