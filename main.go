package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"log"
	database "mtgpoolservice/db"
	"mtgpoolservice/logging"
	"mtgpoolservice/models/entities"
	"mtgpoolservice/models/mtgjson"
	"mtgpoolservice/routers"
	"mtgpoolservice/routers/api"
	"mtgpoolservice/setting"
	"net/http"
	"time"
)

func init() {
	setting.Setup()
	logging.Setup()
	database.Init()
}

func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	// Check for DB Update
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Day().Do(CheckAndUpdateSets)

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()
}

func CheckAndUpdateSets() {
	log.Println("Check MTGJSON remote version")
	version := mtgjson.Version{}

	resp, err := http.Get("http://mtgjson.com/api/v5/Meta.json")
	if err != nil {
		log.Println("Could not fetch the MTGJson version")
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		log.Println("error while unmarshalling version", err)
		return
	}

	log.Println(version)
	entity := mapToVersion(version)
	lastVersion, err := database.FetchLastVersion()
	if err != nil {
		fmt.Print("Could not fetch last version", err)
		return
	}
	if entity.IsNewer(lastVersion) {
		if err := database.GetDB().Save(&entity).Error; err != nil {
			fmt.Printf("could not save the version %w\n", err)
			return
		}

		/**
		Update the DB
		*/
		err := api.UpdateSets()
		if err != nil {
			fmt.Print("main.CheckAndUpdateSets() - ERROR: %w", err)
		}
	}
}

func mapToVersion(version mtgjson.Version) entities.Version {
	date := time.Time(version.Data.Date)
	v := entities.Version{
		Date:            date,
		SemanticVersion: version.Data.Version,
	}
	return v
}
