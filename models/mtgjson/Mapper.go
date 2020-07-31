package mtgjson

import (
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"mtgpoolservice/models/entities"
)

func MapMTGJsonSetToEntity(mtgJsonSet MTGJsonSet) entities.Set {
	s := entities.Set{
		Code:               mtgJsonSet.Code,
		Name:               mtgJsonSet.Name,
		Type:               mtgJsonSet.Type,
		ReleaseDate:        mtgJsonSet.ReleaseDate,
		BaseSetSize:        mtgJsonSet.BaseSetSize,
		Cards:              MakeCards(mtgJsonSet.Code, mtgJsonSet.Cards),
		Sheets:             MakeSheets(mtgJsonSet.Code, mtgJsonSet.Booster.Default.Sheets),
		PackConfigurations: MakePackConfigurations(mtgJsonSet.Booster.Default.Boosters),
	}
	return s
}

func MakePackConfigurations(configurations []PackConfiguration) postgres.Jsonb {
	jsonContent, _ := json.Marshal(configurations)
	return postgres.Jsonb{RawMessage: jsonContent}
}

func MakeSheets(code string, sheets map[string]Sheet) (ret []entities.Sheet) {
	for name, sheet := range sheets {
		sh := entities.Sheet{
			ID:            code + "_" + name,
			SetID:         code,
			Name:          name,
			BalanceColors: sheet.BalanceColors,
			Foil:          sheet.Foil,
			TotalWeight:   sheet.TotalWeight,
			Cards:         MakeSheetCards(code+"_"+name, sheet.Cards),
		}
		ret = append(ret, sh)
	}
	return
}

func MakeSheetCards(sheetId string, cards SheetCards) (ret []entities.SheetCard) {
	for _, card := range cards {
		sc := entities.SheetCard{
			SheetID: sheetId,
			UUID:    card.UUID,
			Weight:  card.Weight,
		}
		ret = append(ret, sc)
	}
	return
}

func MakeCards(code string, cards []Card) (ret []entities.Card) {
	for _, card := range cards {
		mappedCard := entities.Card{
			SetID:             code,
			UUID:              card.UUID,
			Name:              card.Name,
			Number:            card.Number,
			Layout:            card.Layout,
			Loyalty:           card.Loyalty,
			Power:             card.Power,
			Toughness:         card.Toughness,
			ConvertedManaCost: card.ConvertedManaCost,
			Type:              card.Types[0], //TODO: check if always true
			ManaCost:          card.ManaCost,
			Rarity:            card.Rarity,
			Side:              card.Side,
			IsAlternative:     card.IsAlternative,
			Colors:            MakeColors(card.Colors),
			Color:             GetColor(card.Colors),
		}

		ret = append(ret, mappedCard)
	}

	return
}

func MakeColors(colors []string) (ret []entities.Color) {
	for _, color := range colors {
		c := entities.Color{
			ID: color,
		}
		ret = append(ret, c)
	}
	return
}

func GetColor(colors []string) string {
	if len(colors) == 0 {
		return "colorless"
	}
	switch len(colors) {
	case 0:
		return "colorless"
	case 1:
		return colors[0]
	default:
		return "multicolor"
	}
}
