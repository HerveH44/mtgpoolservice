package routers

import (
	"github.com/gin-gonic/gin"
	"mtgpoolservice/routers/api"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", api.Ping)

	//Update DB
	r.GET("/refresh", api.Refresh)
	r.GET("/refresh/:setCode", api.RefreshSet)

	r.POST("/regular", api.RegularPacks)
	r.POST("/cube", api.CubePacks)

	return r
}
