package entities

type Card struct {
	UUID              string `json:"uuid" gorm:"primary_key"`
	SetID             string `json:"setCode"`
	Name              string `json:"name"`
	Number            string `json:"number"`
	Layout            string `json:"layout"`
	Loyalty           string `json:"loyalty"`
	Power             string `json:"power"`
	Toughness         string `json:"toughness"`
	ConvertedManaCost int    `json:"cmc"`
	Type              string `json:"type"`
	ManaCost          string `json:"manaCost"`
	Rarity            string `json:"rarity"`
	Side              string `json:"side"`
	IsAlternative     bool   `json:"isAlternative"`
	Color             string `json:"color"`
	ScryfallID        string `json:"scryfallId"`
	URL               string `json:"url"`
	Cubable           bool   `json:"-" gorm:"index:cubable_idx"`
	FaceName          string `json:"-" gorm:"index:cubable_idx"`
}
