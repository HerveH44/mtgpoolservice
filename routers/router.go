package routers

import (
	"github.com/gin-gonic/gin"
	"mtgpoolservice/routers/api"
)

func InitRouter(regularPacksController api.RegularController, setController api.SetController, importerController api.ImporterController, cubeController api.CubeController, chaosController api.ChaosController) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", api.Ping)

	r.GET("/infos", setController.GetInfos)
	r.GET("/sets", setController.GetAvailableSets)
	r.GET("/sets/latest", setController.GetLatestSet)
	r.GET("/about", setController.GetVersion)

	r.GET("/import", importerController.ImportAllSets)
	r.GET("/import/:setCode", importerController.ImportSet)

	r.POST("/regular", regularPacksController.RegularPacks)
	r.POST("/chaos", chaosController.ChaosPacks)
	r.POST("/cube", cubeController.CubePacks)
	r.POST("/cubelist", cubeController.CubeList)

	return r
}
