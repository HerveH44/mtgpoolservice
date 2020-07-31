package entities

type Color struct {
	ID string `gorm:"primary_key"`
}

type Card struct {
	SetID             string
	UUID              string
	Name              string
	Number            string
	Layout            string
	Loyalty           string
	Power             string
	Toughness         string
	ConvertedManaCost float32
	Type              string
	ManaCost          string
	Rarity            string
	Side              string
	IsAlternative     bool
	Color             string
	Colors            []Color `gorm:"many2many:card_colors;PRELOAD:true"`
	//URL               string
	//Names             []string  `gorm:"many2many:card_names;PRELOAD:true"`
}
