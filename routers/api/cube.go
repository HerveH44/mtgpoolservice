package api

import (
	"github.com/gin-gonic/gin"
	"mtgpoolservice/models"
	"mtgpoolservice/services"
	"net/http"
)

func CubePacks(c *gin.Context) {
	var request models.CubeDraftRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret := make([]models.Pool, 0)
	for p := 0; p < int(request.Players); p++ {
		packs, err := services.MakeCubePacks(request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ret = append(ret, packs...)
	}
	c.JSON(http.StatusOK, ret)
}
