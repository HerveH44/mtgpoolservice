package models

type Sheet struct {
	BalanceColors bool                `json:"balanceColors"`
	Cards         map[MTGCardUUID]int `json:"cards"`
	Foil          bool                `json:"foil"`
	TotalWeight   int                 `json:"totalWeight"`
}

type SheetName = string
type MTGCardUUID = string

type Configuration struct {
	Contents map[SheetName]int `json:"contents"`
	Weight   int               `json:"weight"`
}

type BoosterPackConfiguration struct {
	Boosters            []Configuration
	BoostersTotalWeight int
	Sheets              map[SheetName]Sheet
}

type Set struct {
	Code        string                              `json:"code"`
	Name        string                              `json:"name"`
	Type        string                              `json:"type"`
	ReleaseDate string                              `json:"releaseDate"`
	BaseSetSize int                                 `json:"baseSetSize"`
	Cards       []Card                              `json:"cards"`
	Booster     map[string]BoosterPackConfiguration `json:"booster"`
}

func (set *Set) addCard(card Card) []Card {
	set.Cards = append(set.Cards, card)
	return set.Cards
}

type Card struct {
	UUID              string   `json:"uuid"`
	Name              string   `json:"name"`
	Number            string   `json:"number"`
	Layout            string   `json:"layout"`
	Names             []string `json:"names"`
	Loyalty           string   `json:"loyalty"`
	Power             string   `json:"power"`
	Toughness         string   `json:"toughness"`
	ConvertedManaCost float32  `json:"convertedManaCost"`
	Colors            []string `json:"colors"`
	Types             []string `json:"types"`
	Supertypes        []string `json:"supertypes"`
	ManaCost          string   `json:"manaCost"`
	URL               string   `json:"url"`
	Rarity            string   `json:"rarity"`
	ScryfallID        string   `json:"scryfallId"`
	Side              string   `json:"side"`
	IsAlternative     bool     `json:"isAlternative"`
}
