package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
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
	//db.LogMode(setting.DatabaseSetting.Log)
	//db.SetLogger(logging.GetLogger())

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

func CheckCubeCards(names []string) (missingCardNames []string) {
	faceNames := GetFaceNames(names[:])

	jobs := make(chan string, len(faceNames))
	missingCards := make(chan string)
	var wg sync.WaitGroup

	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go worker(jobs, missingCards, &wg)
	}
	go pushToMissingCards(missingCards, &missingCardNames, &wg)

	for _, name := range faceNames {
		jobs <- name
	}
	close(jobs)
	wg.Wait()
	return
}

func pushToMissingCards(missingCards <-chan string, missingCardNames *[]string, wg *sync.WaitGroup) {
	defer wg.Done()
	for c := range missingCards {
		*missingCardNames = append(*missingCardNames, c)
	}
}

func worker(jobs <-chan string, missingCard chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		_, err := GetCardByFacename(j)
		if err != nil {
			missingCard <- j
		}
	}
}

func GetCardsByName(names []string) (cr []models.CardResponse, missingCards []string) {
	faceNames := GetFaceNames(names[:])
	cards := make([]entities.Card, 0)
	for _, name := range faceNames {
		card, err := GetCardByFacename(name)
		if err != nil {
			log.Println("could not find card with face_name", name)
			missingCards = append(missingCards, name)
		} else {
			cards = append(cards, card)
		}
	}

	for i, _ := range cards {
		cardResponse := models.CardResponse{
			Card: &cards[i],
			Id:   uuid.New().String(),
			Foil: false,
		}
		cr = append(cr, cardResponse)
	}
	return
}

func GetCardByFacename(name string) (card entities.Card, err error) {
	err = DB.Where("cubable = true AND face_name = ?", name).First(&card).Error
	if err != nil {
		log.Println("could not find card with face_name", name)
	}
	return
}

func GetFaceNames(names []string) (faceNames []string) {
	for _, name := range names {
		facename := strings.ToLower(strings.Split(name, " // ")[0])
		faceNames = append(faceNames, facename)
	}
	return
}
