package db

import (
	"fmt"
	"mtgpoolservice/logging"
	"mtgpoolservice/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Opening a database and save the reference to `Database` struct.
func ConnectDB(settings setting.Settings) (db *gorm.DB, err error) {
	db, err = gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		settings.Database.Host,
		settings.Database.Port,
		settings.Database.User,
		settings.Database.Name,
		settings.Database.Password,
		settings.Database.SslMode))

	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(10)
	db.LogMode(settings.Database.Log)
	db.SetLogger(logging.GetLogger())

	db.AutoMigrate(&Set{})
	db.AutoMigrate(&Card{})
	db.AutoMigrate(&Sheet{})
	db.AutoMigrate(&SheetCard{})
	db.AutoMigrate(&Version{})
	return
}
