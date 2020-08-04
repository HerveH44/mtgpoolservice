package entities

type Color struct {
	ID string `gorm:"primary_key"`
}

type Card struct {
	UUID              string  `json:"uuid" gorm:"primary_key"`
	SetID             string  `json:"setCode"`
	Name              string  `json:"name"`
	Number            string  `json:"number"`
	Layout            string  `json:"layout"`
	Loyalty           string  `json:"loyalty"`
	Power             string  `json:"power"`
	Toughness         string  `json:"toughness"`
	ConvertedManaCost float32 `json:"convertedManaCost"`
	Type              string  `json:"type"`
	ManaCost          string  `json:"manaCost"`
	Rarity            string  `json:"rarity"`
	Side              string  `json:"side"`
	IsAlternative     bool    `json:"isAlternative"`
	Color             string  `json:"color"`
	Colors            []Color `json:"colors" gorm:"many2many:card_colors;PRELOAD:true"`
	ScryfallID        string  `json:"scryfallId"`
	URL               string  `json:"url"`
	Cubable           bool    `json:"-"`
	//Names             []string  `gorm:"many2many:card_names;PRELOAD:true"`
}
