package api

import (
	"github.com/gin-gonic/gin"
	"mtgpoolservice/logging"
	"mtgpoolservice/models"
	"mtgpoolservice/services"
	"net/http"
)

func RegularPacks(c *gin.Context) {
	var request models.RegularRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logging.Warn(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret := make([]models.CardPool, 0)
	for p := 0; p < request.Players; p++ {
		packs, err := services.MakePacks(request.Sets)
		if err != nil {
			logging.Warn(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ret = append(ret, packs...)
	}
	c.JSON(http.StatusOK, ret)
}
