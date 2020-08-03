package routers

import (
	"github.com/gin-gonic/gin"
	"mtgpoolservice/routers/api"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", api.Ping)

	//Update DB
	r.GET("/refresh", api.RefreshSetsInDB)
	r.GET("/refresh/:setCode", api.RefreshSet)

	r.POST("/regular", api.RegularPacks)
	r.POST("/cube", api.CubePacks)

	return r
}
