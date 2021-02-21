package db

import (
	"fmt"
	"mtgpoolservice/setting"

	"gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(settings setting.Settings) (db *gorm.DB, err error) {
	db, err = gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		settings.Database.Host,
		settings.Database.Port,
		settings.Database.User,
		settings.Database.Name,
		settings.Database.Password,
		settings.Database.SslMode)), &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           nil,
		FullSaveAssociations:                     false,
		Logger:                                   logger.Default.LogMode(logger.Info),
		NowFunc:                                  nil,
		DryRun:                                   false,
		PrepareStmt:                              false,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 false,
		AllowGlobalUpdate:                        false,
		QueryFields:                              false,
		CreateBatchSize:                          0,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Set{})
	db.AutoMigrate(&Card{})
	db.AutoMigrate(&Sheet{})
	db.AutoMigrate(&SheetCard{})
	db.AutoMigrate(&Version{})
	return
}
