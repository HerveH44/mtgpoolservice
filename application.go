package main

import (
	"fmt"
	database "mtgpoolservice/db"
	"mtgpoolservice/importer"
	"mtgpoolservice/importer/mtgjson"
	"mtgpoolservice/logging"
	"mtgpoolservice/pool"
	"mtgpoolservice/routers"
	"mtgpoolservice/routers/api"
	"mtgpoolservice/setting"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	settings := setting.GetSettings()
	logging.Setup(settings.App)

	db, err := database.ConnectDB(settings)
	if err != nil {
		log.Fatal("Could not initialize DB ", err)
	}

	setRepository := database.NewSetRepository(db)
	cardRepository := database.NewCardRepository(db)
	versionRepository := database.NewVersionRepository(db)

	packService := pool.NewPackService(setRepository, cardRepository)
	mtgJsonService := mtgjson.NewMTGJsonService(settings.MTGJsonEndpoint)
	importerFacade := importer.NewImporterFacade(mtgJsonService, setRepository, versionRepository)

	regularPackController := api.NewRegularController(packService)
	setController := api.NewSetController(setRepository, versionRepository)
	importerController := api.NewImporterController(importerFacade)
	cubeController := api.NewCubeController(packService)
	chaosController := api.NewChaosController(packService)

	routersInit := routers.InitRouter(regularPackController, setController, importerController, cubeController, chaosController)

	gin.SetMode(settings.Server.RunMode)
	readTimeout := settings.Server.ReadTimeout
	writeTimeout := settings.Server.WriteTimeout
	port := settings.Server.HttpPort
	endPoint := fmt.Sprintf(":%d", port)
	maxHeaderBytes := 1 << 20

	// Check for DB Update
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Day().Do(func() {
		if updateError := importerFacade.UpdateSets(false); updateError != nil {
			log.Error("Could not update sets. Error", updateError)
		}
	})

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Info("start http server listening", endPoint)

	server.ListenAndServe()
}
