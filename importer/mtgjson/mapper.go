package mtgjson

import (
	"encoding/json"
	"fmt"
	"mtgpoolservice/db"
	"mtgpoolservice/utils"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

func MapMTGJsonVersionToVersion(version Meta) db.Version {
	date := time.Time(version.Date)
	v := db.Version{
		Date:            date,
		SemanticVersion: version.Version,
	}
	return v
}

func MapMTGJsonSetToEntity(mtgJsonSet *MTGJsonSet, isCubable func(string, *Card) bool) *db.Set {
	set := db.Set{
		Code:               mtgJsonSet.Code,
		Name:               mtgJsonSet.Name,
		Type:               mtgJsonSet.Type,
		ReleaseDate:        time.Time(mtgJsonSet.ReleaseDate),
		BaseSetSize:        mtgJsonSet.BaseSetSize,
		Cards:              MakeCards(mtgJsonSet.Code, mtgJsonSet.Cards, isCubable),
		Sheets:             MakeSheets(mtgJsonSet.Code, mtgJsonSet.Booster.Default.Sheets),
		PackConfigurations: MakePackConfigurations(mtgJsonSet.Booster.Default.Boosters),
	}
	return &set
}

func MakePackConfigurations(configurations []PackConfiguration) postgres.Jsonb {
	jsonContent, _ := json.Marshal(configurations)
	return postgres.Jsonb{RawMessage: jsonContent}
}

func MakeSheets(code string, sheets map[string]Sheet) (ret []db.Sheet) {
	for name, sheet := range sheets {
		sh := db.Sheet{
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

func MakeSheetCards(sheetId string, cards SheetCards) (ret []db.SheetCard) {
	for _, card := range cards {
		sc := db.SheetCard{
			SheetID: sheetId,
			UUID:    card.UUID,
			Weight:  card.Weight,
		}
		ret = append(ret, sc)
	}
	return
}

func MakeCards(code string, cards []Card, isCubable func(string, *Card) bool) (ret []db.Card) {
	for _, card := range cards {
		if card.IsPromo || card.IsAlternative {
			continue
		}

		props := getDoubleFacedProps(&card, cards)

		mappedCard := db.Card{
			UUID:              card.UUID,
			Name:              card.Name,
			FaceName:          MakeFaceName(card.FaceName, card.Name),
			Color:             GetColor(card.Colors),
			SetID:             code,
			ConvertedManaCost: int(card.ConvertedManaCost),
			Number:            card.Number,
			Type:              card.Types[0], //TODO: check if always true
			ManaCost:          card.ManaCost,
			Rarity:            GetRarity(&card),
			URL:               fmt.Sprintf("https://api.scryfall.com/cards/%s?format=image", card.Identifiers.ScryfallID),
			ScryfallID:        db.Identifier(card.Identifiers.ScryfallID),
			Layout:            card.Layout,
			IsDoubleFaced:     props.IsDoubleFace,
			FlippedCardURL:    props.flippedCardURL,
			FlippedIsBack:     props.flippedIsBack,
			FlippedNumber:     props.flippedNumber,
			Text:              card.Text,
			Loyalty:           card.Loyalty,
			Power:             card.Power,
			Toughness:         card.Toughness,
			Side:              card.Side,
			IsAlternative:     card.IsAlternative,
			Cubable:           isCubable(code, &card),
		}

		ret = append(ret, mappedCard)
	}

	return
}

type doubleFacedProps struct {
	IsDoubleFace   bool
	flippedCardURL string
	flippedIsBack  bool
	flippedNumber  string
}

func getDoubleFacedProps(c *Card, cards []Card) (props doubleFacedProps) {
	props.IsDoubleFace = regexp.MustCompile("/^modal_dfc$|^double-faced$|^transform$|^flip$|^meld$|/").MatchString(c.Layout)
	if !props.IsDoubleFace {
		return
	}
	names := strings.Split(c.Name, " // ")
	if len(names) < 2 {
		return
	}
	for _, card := range cards {
		if names[1] != card.FaceName {
			continue
		}

		scryfallId := card.Identifiers.ScryfallID
		props.flippedCardURL = fmt.Sprintf("https://api.scryfall.com/cards/%s?format=image", scryfallId)
		if regexp.MustCompile("/^modal_dfc$|^double-faced$|^transform$|/").MatchString(c.Layout) {
			props.flippedCardURL += "&face=back"
			props.flippedNumber = card.Number
			props.flippedIsBack = true
		}

		if regexp.MustCompile("^meld$").MatchString(c.Layout) {
			props.flippedNumber = card.Number
		}

		break
	}

	return
}

func GetRarity(c *Card) string {
	if utils.Include(c.Supertypes, "Basic") {
		return "Basic"
	}
	return strings.Title(c.Rarity)
}

func MakeFaceName(faceName string, name string) string {
	if faceName != "" {
		return strings.ToLower(faceName)
	}
	return NameToFaceName(name)
}

func NameToFaceName(name string) string {
	return strings.ToLower(strings.Split(name, " // ")[0])
}

func GetColor(colors []string) string {
	if len(colors) == 0 {
		return "Colorless"
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
