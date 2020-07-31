package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"mtgpoolservice/models"
	"mtgpoolservice/models/entities"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {
	dsn := "host=localhost port=5432 user=postgres dbname=mtgpoolservice sslmode=disable"
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Println("entities err: ", err)
	}
	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)

	db.AutoMigrate(&entities.Set{})
	db.AutoMigrate(&entities.Card{})
	db.AutoMigrate(&entities.Color{})
	db.AutoMigrate(&entities.BoosterRule{})
	db.AutoMigrate(&entities.Sheet{})
	db.AutoMigrate(&entities.SheetCard{})
	db.AutoMigrate(&entities.PackConfiguration{})
	db.AutoMigrate(&entities.ConfigurationContent{})
	DB = db
	return DB
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}

func GetSet(setCode string) (s entities.Set, err error) {
	fmt.Printf("fetching set %s\n", setCode)
	err = DB.Where(" code = ?", setCode).Set("gorm:auto_preload", true).First(&s).Error
	return
}

func GetCards(protoCards []entities.ProtoCard) (cr []models.CardResponse, err error) {
	c := make([]entities.Card, 0)
	fmt.Printf("fetching cards %s\n", protoCards)

	uuids := make([]string, 0)
	for _, protoCard := range protoCards {
		uuids = append(uuids, protoCard.UUID)
	}

	err = DB.Where(" uuid IN (?)", uuids).Set("gorm:auto_preload", true).Find(&c).Error

	for i, card := range protoCards {
		cr = append(cr, models.CardResponse{
			Card: *getCardFromSlice(card.UUID, c),
			Id:   uuid.New().String(),
			Foil: protoCards[i].Foil,
		})
	}
	return
}

func getCardFromSlice(uuid string, cards []entities.Card) *entities.Card {
	for _, card := range cards {
		if card.UUID == uuid {
			return &card
		}
	}
	return nil
}

func GetCardsByName(names []string) (cr []models.CardResponse, err error) {
	for _, name := range names {
		var card entities.Card
		err = DB.Where("name ILIKE ?", name).Set("gorm:auto_preload", true).First(&card).Error
		if err != nil {
			return nil, fmt.Errorf("could not find cardResponse with name %s", name)
		}
		cardResponse := models.CardResponse{
			Card: card,
			Id:   uuid.New().String(),
			Foil: false,
		}
		cr = append(cr, cardResponse)
	}

	return
}
