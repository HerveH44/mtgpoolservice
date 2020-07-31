package mtgjson

type Meta struct {
	Date    string `json:"date"`
	Version string `json:"version"`
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
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Block        string `json:"block"`
	ReleaseDate  string `json:"releaseDate"`
	BaseSetSize  int    `json:"baseSetSize"`
	TotalSetSize int    `json:"totalSetSize"`
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
	Layout        string        `json:"layout"`
	Number        string        `json:"number"`
	Power         string        `json:"power,omitempty"`
	Toughness     string        `json:"toughness,omitempty"`
	Printings     []string      `json:"printings"`
	Rarity        string        `json:"rarity"`
	Text          string        `json:"text,omitempty"`
	Type          string        `json:"type"`
	Types         []string      `json:"types"`
	Subtypes      []string      `json:"subtypes"`
	Supertypes    []interface{} `json:"supertypes"`
	IsReprint     bool          `json:"isReprint,omitempty"`
	IsAlternative bool          `json:"isAlternative,omitempty"`
	FrameEffects  []string      `json:"frameEffects,omitempty"`
	OtherFaceIds  []string      `json:"otherFaceIds,omitempty"`
	Side          string        `json:"side,omitempty"`
	Variations    []string      `json:"variations,omitempty"`
}
