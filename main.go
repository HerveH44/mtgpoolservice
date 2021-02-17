package main

import (
	"fmt"
	"log"
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

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
)

func main() {
	var err error
	settings := setting.GetSettings()
	err = logging.Setup(settings.App)
	if err != nil {
		log.Fatalf("Could not initialize logger %s", err)
	}

	db, err := database.ConnectDB(settings)
	if err != nil {
		logging.Fatal("Could not initialize DB %s", err)
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
			logging.Error("Could not update sets. Error", updateError)
		}
	})

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	logging.Info("[info] start http server listening", endPoint)

	server.ListenAndServe()
}
