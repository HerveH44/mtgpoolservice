package mtgjson

import (
	"encoding/json"
)

type Set struct {
	Code        string       `json:"code" gorm:"primary_key"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	ReleaseDate string       `json:"releaseDate"`
	BaseSetSize int          `json:"baseSetSize"`
	Cards       []Card       `json:"cards" gorm:"foreignkey:SetID;PRELOAD:true"`
	Booster     BoosterRules `json:"booster" gorm:"foreignkey:SetID;PRELOAD:true"`
}

type Color struct {
	ID string `gorm:"primary_key" json:"-,"`
}

func (i *Color) UnmarshalJSON(data []byte) error {
	s := string(data)
	i.ID = s[1 : len(s)-1]
	return nil
}

func (i *Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ID)
}

type Type struct {
	ID string `gorm:"primary_key"`
}

func (i *Type) UnmarshalJSON(data []byte) error {
	s := string(data)
	i.ID = s[1 : len(s)-1]
	return nil
}

func (i *Type) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ID)
}

type Supertype struct {
	ID string `gorm:"primary_key"`
}

func (i *Supertype) UnmarshalJSON(data []byte) error {
	s := string(data)
	i.ID = s[1 : len(s)-1]
	return nil
}

func (i *Supertype) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ID)
}

type Name struct {
	ID string `gorm:"primary_key"`
}

func (n *Name) UnmarshalJSON(data []byte) error {
	s := string(data)
	n.ID = s[1 : len(s)-1]
	return nil
}

func (i *Name) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.ID)
}

type Card struct {
	SetID             string
	UUID              string  `json:"uuid" gorm:"primary_key"`
	Name              string  `json:"name"`
	Number            string  `json:"number"`
	Layout            string  `json:"layout"`
	Loyalty           string  `json:"loyalty"`
	Power             string  `json:"power"`
	Toughness         string  `json:"toughness"`
	ConvertedManaCost float32 `json:"convertedManaCost"`
	Type              string  `json:"type`
	ManaCost          string  `json:"manaCost"`
	URL               string  `json:"url"`
	Rarity            string  `json:"rarity"`
	Side              string  `json:"side"`
	IsAlternative     bool    `json:"isAlternative"`
	Colors            []Color `json:"colors" gorm:"many2many:card_colors;PRELOAD:true"`
	//Names             []Name      `json:"names" gorm:"many2many:card_names;PRELOAD:true"`
}
