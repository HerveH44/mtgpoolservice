package api

import (
	"github.com/gin-gonic/gin"
	"mtgpoolservice/models"
	"mtgpoolservice/services"
	"net/http"
)

func CubePacks(c *gin.Context) {
	var req models.CubeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if int(req.PlayerPackSize)*int(req.Players)*int(req.Packs) > len(req.Cubelist) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cube list is too small"})
		return
	}

	packs, err := services.MakeCubePacks(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packs)
}

func CubeList(c *gin.Context) {
	var request models.CubeListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errs := services.CheckCubeList(request)
	if len(errs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": errs})
	} else {
		c.Status(200)
	}
}
