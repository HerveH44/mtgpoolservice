package routers

import (
	"github.com/gin-gonic/gin"
	"mtgpoolservice/routers/api"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", api.Ping)
	r.GET("/about", api.Version)

	//Update DB
	r.GET("/import", api.ImportAllSets)
	r.GET("/import/:setCode", api.ImportSet)

	r.GET("/sets", api.AvailableSets)
	r.GET("/sets/latest", api.LatestSet)
	r.POST("/regular", api.RegularPacks)
	r.POST("/cube", api.CubePacks)
	r.POST("/cubelist", api.CubeList)
	r.POST("/chaos", api.ChaosPacks)

	return r
}
