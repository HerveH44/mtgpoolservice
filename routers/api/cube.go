package api

import (
	"github.com/gin-gonic/gin"
	"mtgpoolservice/models"
	"mtgpoolservice/services"
	"net/http"
)

func CubePacks(c *gin.Context) {
	var request models.CubeSealedRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret := make([]models.Pool, request.Players)
	for p := 0; p < int(request.Players); p++ {
		pa := make(models.Pool, 0)
		packs, err := services.MakeCubePacks(request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		for _, pack := range packs {
			pa = append(pa, pack...)
		}
		ret[p] = pa
	}
	c.JSON(http.StatusOK, ret)
}
