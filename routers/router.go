package routers

import (
	"bytes"
	"io"
	"io/ioutil"
	"mtgpoolservice/routers/api"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InitRouter(regularPacksController api.RegularController, setController api.SetController, importerController api.ImporterController, cubeController api.CubeController, chaosController api.ChaosController) *gin.Engine {
	r := gin.Default()
	r.Use(RequestLoggerMiddleware())

	r.GET("/ping", api.Ping)

	r.GET("/infos", setController.GetInfos)
	r.GET("/sets", setController.GetAvailableSets)
	r.GET("/sets/latest", setController.GetLatestSet)
	r.GET("/about", setController.GetVersion)

	r.GET("/import", importerController.ImportAllSets)

	r.POST("/regular", regularPacksController.RegularPacks)
	r.POST("/chaos", chaosController.ChaosPacks)
	r.POST("/cube", cubeController.CubePacks)
	r.POST("/cubelist", cubeController.CubeList)

	return r
}

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		log.Debug(string(body))
		log.Debug(c.Request.Header)
		c.Next()
	}
}
