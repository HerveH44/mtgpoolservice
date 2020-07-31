package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"log"
	database "mtgpoolservice/db"
	"mtgpoolservice/models"
	"mtgpoolservice/models/entities"
	"mtgpoolservice/models/mtgjson"
	"mtgpoolservice/services"
	"net/http"
)

var db = database.Init()

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/refresh", func(context *gin.Context) {
		log.Println("Refresh data")
		resp, err := http.Get("http://mtgjson.com/api/v5/AllPrintings.json")
		if err != nil {
			log.Println("Could not fetch the MTGJson allPrintings")
			context.JSON(500, "unexpected error")
			return
		}
		defer resp.Body.Close()

		allPrintings := new(mtgjson.AllPrintings)
		if err := json.NewDecoder(resp.Body).Decode(allPrintings); err != nil {
			log.Println("error while unmarshalling allPrintings", err)
		}
		var setsNumber = len(allPrintings.Data)
		log.Println("sets found", setsNumber)

		i := 0
		for setName, set := range allPrintings.Data {
			i++
			fmt.Printf("%d/%d - saving set %s\n", i, setsNumber, setName)
			if err := db.Save(&set).Error; err != nil {
				fmt.Printf("could not save the card %s - %s\n", setName, err)
			}
		}
	})
	r.GET("/test", func(context *gin.Context) {
		log.Println("Refresh data")
		resp, err := http.Get("http://mtgjson.com/api/v5/ISD.json")
		if err != nil {
			log.Println("Could not fetch the MTGJson monoSet")
			context.JSON(500, "unexpected error: Could not fetch the MTGJson monoSet")
			return
		}
		defer resp.Body.Close()

		monoSet := new(mtgjson.MonoSet)
		if err := json.NewDecoder(resp.Body).Decode(monoSet); err != nil {
			log.Println("main: error while unmarshalling monoSet", err)
		}

		log.Println("saving cards from ", monoSet.Data.Name)
		entity := MapMTGJsonSetToEntity(monoSet.Data)
		if err := db.Save(&entity).Error; err != nil {
			log.Fatal("main: could not save the card", monoSet.Data.Name, err)
		}
	})
	r.POST("/regular/draft", func(c *gin.Context) {
		var request models.RegularRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ret := make([][]models.Pool, 0)
		for p := 0; p < request.Players; p++ {
			packs, err := services.MakePacks(request.Sets)
			if err != nil {
				c.JSON(http.StatusBadRequest, err)
			}
			ret = append(ret, packs)
		}
		c.JSON(http.StatusOK, ret)
	})
	r.POST("/regular/sealed", func(c *gin.Context) {
		var request models.RegularRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ret := make([]models.Pool, request.Players)
		for p := 0; p < request.Players; p++ {
			pa := make(models.Pool, 0)
			packs, err := services.MakePacks(request.Sets)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			for _, pack := range packs {
				pa = append(pa, pack...)
			}
			ret[p] = pa
		}
		c.JSON(http.StatusOK, ret)
	})
	r.POST("/cube/sealed", func(c *gin.Context) {
		var request models.CubeSealedRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ret := make([]models.Pool, request.Players)
		for p := 0; p < int(request.Players); p++ {
			pa := make(models.Pool, 0)
			packs, err := services.MakeCubePacks(request)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			for _, pack := range packs {
				pa = append(pa, pack...)
			}
			ret[p] = pa
		}
		c.JSON(http.StatusOK, ret)
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func MapMTGJsonSetToEntity(mtgJsonSet mtgjson.MTGJsonSet) entities.Set {
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

func MakePackConfigurations(configurations []mtgjson.PackConfiguration) postgres.Jsonb {
	jsonContent, _ := json.Marshal(configurations)
	return postgres.Jsonb{RawMessage: jsonContent}
}

func MakeSheets(code string, sheets map[string]mtgjson.Sheet) (ret []entities.Sheet) {
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

func MakeSheetCards(sheetId string, cards mtgjson.SheetCards) (ret []entities.SheetCard) {
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

func MakeCards(code string, cards []mtgjson.Card) (ret []entities.Card) {
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
