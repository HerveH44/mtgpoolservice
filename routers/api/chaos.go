package api

import (
	"github.com/gin-gonic/gin"
	"mtgpoolservice/models"
	"mtgpoolservice/services"
	"net/http"
)

func ChaosPacks(c *gin.Context) {
	var req models.ChaosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packs, err := services.MakeChaosPacks(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packs)

}
