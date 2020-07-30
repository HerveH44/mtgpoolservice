package common

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"mtgpoolservice/models"
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
		fmt.Println("db err: ", err)
	}
	db.DB().SetMaxIdleConns(10)
	db.LogMode(true)

	db.AutoMigrate(&models.Set{})
	db.AutoMigrate(&models.Card{})
	db.AutoMigrate(&models.Color{})
	db.AutoMigrate(&models.Type{})
	db.AutoMigrate(&models.Supertype{})
	db.AutoMigrate(&models.BoosterRule{})
	db.AutoMigrate(&models.PackConfiguration{})
	db.AutoMigrate(&models.Sheet{})
	db.AutoMigrate(&models.SheetCard{})
	db.AutoMigrate(&models.ConfigurationContent{})
	DB = db
	return DB
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}

func GetSet(setCode string) (s models.Set, err error) {
	fmt.Printf("fetching set %s\n", setCode)
	err = DB.Where(" code = ?", setCode).Set("gorm:auto_preload", true).First(&s).Error
	return
}

func GetCards(protoCards []models.ProtoCard) (cr []models.CardResponse, err error) {
	c := make([]models.Card, 0)
	fmt.Printf("fetching cards %s\n", protoCards)

	uuids := make([]string, 0)
	for _, protoCard := range protoCards {
		uuids = append(uuids, protoCard.UUID)
	}

	err = DB.Where(" uuid IN (?)", uuids).Set("gorm:auto_preload", true).Find(&c).Error

	for i, card := range c {
		cr = append(cr, models.CardResponse{
			Card: card,
			Id:   uuid.New().String(),
			Foil: protoCards[i].Foil,
		})
	}
	return
}
