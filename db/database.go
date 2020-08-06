package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"mtgpoolservice/logging"
	"mtgpoolservice/models"
	"mtgpoolservice/models/entities"
	"mtgpoolservice/setting"
	"strings"
	"sync"
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
	db.SetLogger(logging.GetLogger())

	db.AutoMigrate(&entities.Set{})
	db.AutoMigrate(&entities.Card{})
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

func getSets() (*[]entities.Set, error) {
	s := make([]entities.Set, 0)
	if err := DB.Find(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func fetchSet(setCode string) (*entities.Set, error) {
	var s entities.Set
	err := DB.Where(" code = ?", setCode).Set("gorm:auto_preload", true).First(&s).Error
	return &s, err
}

func FetchLastVersion() (*entities.Version, error) {
	var v entities.Version
	err := DB.Order("date DESC").First(&v).Error
	return &v, err
}

func addToCardPool(cards chan entities.Card, cr *models.CardPool, wg *sync.WaitGroup) {
	defer wg.Done()
	for c := range cards {
		cr.Add(&c, false)
	}
}

func addToMissingCards(missingCards <-chan string, missingCardNames *[]string, wg *sync.WaitGroup) {
	defer wg.Done()
	for c := range missingCards {
		*missingCardNames = append(*missingCardNames, c)
	}
}

func getCard(jobs <-chan string, missingCards chan<- string, foundCards chan<- entities.Card, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		c, err := GetCardWithName(j)
		if err != nil {
			missingCards <- j
		} else {
			foundCards <- c
		}
	}
}

func GetCardsByName(names []string) (cr models.CardPool, missingCardNames []string) {
	jobs := make(chan string, len(names))
	missingCards := make(chan string)
	cards := make(chan entities.Card)
	var wg sync.WaitGroup

	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go getCard(jobs, missingCards, cards, &wg)
	}

	go addToMissingCards(missingCards, &missingCardNames, &wg)
	go addToCardPool(cards, &cr, &wg)

	for _, name := range names {
		jobs <- name
	}

	close(jobs)
	wg.Wait()
	return
}

func getCardWithName(name string) (card entities.Card, err error) {
	if isMultiCard := strings.ContainsAny(name, "/"); isMultiCard {
		err = DB.Where("cubable = true AND name ILIKE ?", name).First(&card).Error
	} else {
		err = DB.Where("cubable = true AND face_name = ?", name).First(&card).Error
	}
	if err != nil {
		log.Println("could not find card with name", name)
	}
	return
}
