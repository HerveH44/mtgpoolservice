package entities

type Card struct {
	UUID              string `json:"uuid" gorm:"primary_key"`
	Name              string `json:"name"`
	FaceName          string `json:"-" gorm:"index:cubable_idx"`
	Color             string `json:"color"`
	SetID             string `json:"setCode" gorm:"index:pack_idx"`
	ConvertedManaCost int    `json:"cmc"`
	Number            string `json:"number"`
	Type              string `json:"type"`
	ManaCost          string `json:"manaCost"`
	Rarity            string `json:"rarity" gorm:"index:pack_idx"`
	URL               string `json:"url"`
	ScryfallID        string `json:"scryfallId"` // Must be rename as identifiers with inner scryfallId
	Layout            string `json:"layout"`
	// isDoubleFaced
	// flippedCardURL
	// flippedIsBack
	// flippedNumber
	// text
	Loyalty       string `json:"loyalty"`
	Power         string `json:"power"`
	Toughness     string `json:"toughness"`
	Side          string `json:"side"`          // no use ?
	IsAlternative bool   `json:"isAlternative"` // no use ?
	Cubable       bool   `json:"-" gorm:"index:cubable_idx"`
}
