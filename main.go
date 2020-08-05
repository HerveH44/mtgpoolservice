package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"log"
	database "mtgpoolservice/db"
	"mtgpoolservice/logging"
	"mtgpoolservice/routers"
	"mtgpoolservice/services"
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
	port := setting.ServerSetting.HttpPort
	endPoint := fmt.Sprintf(":%d", port)
	maxHeaderBytes := 1 << 20

	// Check for DB Update
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Day().Do(services.CheckAndUpdateSets)

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
