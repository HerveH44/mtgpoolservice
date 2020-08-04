package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"mtgpoolservice/models"
	"mtgpoolservice/models/entities"
	"mtgpoolservice/setting"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {
	db, err := gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Port,
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Name,
		setting.DatabaseSetting.Password))

	if err != nil {
		fmt.Println("entities err: ", err)
	}

	db.DB().SetMaxIdleConns(10)
	db.LogMode(setting.DatabaseSetting.Log)

	db.AutoMigrate(&entities.Set{})
	db.AutoMigrate(&entities.Card{})
	db.AutoMigrate(&entities.Color{})
	db.AutoMigrate(&entities.Sheet{})
	db.AutoMigrate(&entities.SheetCard{})
	db.AutoMigrate(&entities.Version{})
	DB = db
	return DB
}

// Using this function to get a connection, you can create your connection pool here.
func GetDB() *gorm.DB {
	return DB
}

func GetSets() ([]entities.Set, error) {
	s := make([]entities.Set, 0)
	if err := DB.Order("release_date DESC").Find(&s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func fetchSet(setCode string) (*entities.Set, error) {
	fmt.Printf("fetching set %s\n", setCode)
	var s entities.Set
	err := DB.Where(" code = ?", setCode).Set("gorm:auto_preload", true).First(&s).Error
	return &s, err
}

func FetchLastVersion() (*entities.Version, error) {
	var v entities.Version
	err := DB.Order("date DESC").First(&v).Error
	return &v, err
}

func GetCardsByName(names []string) (cr []models.CardResponse, err error) {
	for _, name := range names {
		var card entities.Card
		err = DB.Where("cubable = true AND name ILIKE ?", fmt.Sprint(name, "%s")).Set("gorm:auto_preload", true).First(&card).Error
		if err != nil {
			return nil, fmt.Errorf("could not find card with name like %s", name)
		}
		cardResponse := models.CardResponse{
			Card: &card,
			Id:   uuid.New().String(),
			Foil: false,
		}
		cr = append(cr, cardResponse)
	}

	return
}
