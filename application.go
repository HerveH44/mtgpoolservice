package main

import (
	"fmt"
	database "mtgpoolservice/db"
	"mtgpoolservice/importer"
	"mtgpoolservice/importer/mtgjson"
	"mtgpoolservice/logging"
	"mtgpoolservice/pool"
	"mtgpoolservice/routers"
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
	router := routers.InitRouter()

	importer.
		NewImporterController(importerFacade).
		Register(router.Group("/import", gin.BasicAuth(gin.Accounts{"admin": settings.AdminPassword})))
	database.
		NewDBController(setRepository, versionRepository).
		Register(router.Group("/infos"))
	pool.
		NewPoolController(packService).
		Register(router.Group("/pool"))

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
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Info("start http server listening", endPoint)

	server.ListenAndServe()
}
