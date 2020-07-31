package mtgjson

type MTGJsonSet struct {
	Data struct {
		BaseSetSize int    `json:"baseSetSize"`
		Block       string `json:"block"`
		Booster     struct {
			Default BoosterRule `json:"default"`
		} `json:"booster"`
		Cards []struct {
			UUID              string   `json:"uuid"`
			Name              string   `json:"name"`
			FaceName          string   `json:"faceName,omitempty"`
			ColorIdentity     []string `json:"colorIdentity"`
			Colors            []string `json:"colors"`
			ColorIndicator    []string `json:"colorIndicator,omitempty"`
			ManaCost          string   `json:"manaCost,omitempty"`
			ConvertedManaCost float64  `json:"convertedManaCost"`
			FrameVersion      string   `json:"frameVersion"`
			Identifiers       struct {
				ScryfallID             string `json:"scryfallId"`
				ScryfallIllustrationID string `json:"scryfallIllustrationId"`
				ScryfallOracleID       string `json:"scryfallOracleId"`
			} `json:"identifiers,omitempty"`
			Layout       string        `json:"layout"`
			Number       string        `json:"number"`
			Power        string        `json:"power,omitempty"`
			Toughness    string        `json:"toughness,omitempty"`
			Printings    []string      `json:"printings"`
			Rarity       string        `json:"rarity"`
			Text         string        `json:"text,omitempty"`
			Type         string        `json:"type"`
			Types        []string      `json:"types"`
			Subtypes     []string      `json:"subtypes"`
			Supertypes   []interface{} `json:"supertypes"`
			IsReprint    bool          `json:"isReprint,omitempty"`
			FrameEffects []string      `json:"frameEffects,omitempty"`
			OtherFaceIds []string      `json:"otherFaceIds,omitempty"`
			Side         string        `json:"side,omitempty"`
			Variations   []string      `json:"variations,omitempty"`
		} `json:"cards"`
		Code         string `json:"code"`
		Name         string `json:"name"`
		ReleaseDate  string `json:"releaseDate"`
		TotalSetSize int    `json:"totalSetSize"`
		Type         string `json:"type"`
	} `json:"data"`
	Meta struct {
		Date    string `json:"date"`
		Version string `json:"version"`
	} `json:"meta"`
}
