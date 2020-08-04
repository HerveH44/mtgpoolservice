package mtgjson

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm/dialects/postgres"
	"mtgpoolservice/models/entities"
	"strings"
	"time"
)

func MapMTGJsonSetToEntity(mtgJsonSet MTGJsonSet, isCubable func(string, []string) bool) entities.Set {
	s := entities.Set{
		Code:               mtgJsonSet.Code,
		Name:               mtgJsonSet.Name,
		Type:               mtgJsonSet.Type,
		ReleaseDate:        time.Time(mtgJsonSet.ReleaseDate),
		BaseSetSize:        mtgJsonSet.BaseSetSize,
		Cards:              MakeCards(mtgJsonSet.Code, mtgJsonSet.Cards, isCubable),
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

func MakeCards(code string, cards []Card, isCubable func(string, []string) bool) (ret []entities.Card) {
	for _, card := range cards {
		if card.IsPromo || card.IsAlternative || card.IsStarter {
			continue
		}

		mappedCard := entities.Card{
			SetID:             code,
			UUID:              card.UUID,
			Name:              card.Name,
			Number:            card.Number,
			Layout:            card.Layout,
			Loyalty:           card.Loyalty,
			Power:             card.Power,
			Toughness:         card.Toughness,
			ConvertedManaCost: int(card.ConvertedManaCost),
			Type:              card.Types[0], //TODO: check if always true
			ManaCost:          card.ManaCost,
			Rarity:            strings.Title(card.Rarity),
			Side:              card.Side,
			IsAlternative:     card.IsAlternative,
			Color:             GetColor(card.Colors),
			ScryfallID:        card.Identifiers.ScryfallID,
			URL:               fmt.Sprintf("https://api.scryfall.com/cards/%s?format=image", card.Identifiers.ScryfallID),
			Cubable:           isCubable(code, card.Printings),
			FaceName:          MakeFaceName(card.FaceName, card.Name),
		}

		ret = append(ret, mappedCard)
	}

	return
}

func MakeFaceName(faceName string, name string) string {
	if faceName != "" {
		return faceName
	}
	return strings.ToLower(strings.Split(name, " // ")[0])
}

func GetColor(colors []string) string {
	if len(colors) == 0 {
		return "colorless"
	}
	switch len(colors) {
	case 0:
		return "Colorless"
	case 1:
		return translateColor(colors[0])
	default:
		return "Multicolor"
	}
}

func translateColor(color string) string {
	switch color {
	case "W":
		return "White"
	case "B":
		return "Black"
	case "R":
		return "Red"
	case "U":
		return "Blue"
	case "G":
		return "Green"
	default:
		return "Colorless"
	}
}
