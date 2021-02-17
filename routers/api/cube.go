package api

import (
	"fmt"
	pack "mtgpoolservice/pool"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CubeController interface {
	CubePacks(c *gin.Context)
	CubeList(c *gin.Context)
}

type cubeController struct {
	packService pack.Service
}

func NewCubeController(service pack.Service) CubeController {
	return &cubeController{service}
}

func (cc *cubeController) CubePacks(c *gin.Context) {
	var req CubeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if int(req.PlayerPackSize)*int(req.Players)*int(req.Packs) > len(req.Cubelist) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cube list is too small"})
		return
	}

	packs, err := cc.packService.MakeCubePacks(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packs)
}

func (cc *cubeController) CubeList(c *gin.Context) {
	var request CubeListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errs := cc.packService.CheckCubeList(request)
	if len(errs) > 0 {
		var errorMsg string
		if len(errs) > 10 {
			errorMsg = fmt.Sprintf("Invalid cards: %v and %d more", errs[:10], len(errs)-10)
		} else {
			errorMsg = fmt.Sprintf("Invalid cards: %v", errs)
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
	} else {
		c.Status(200)
	}
}

type CubeRequest struct {
	Cubelist       []string `json:"list"`
	Players        uint     `json:"players"`
	PlayerPackSize uint     `json:"playerPackSize"`
	Packs          uint     `json:"pool"`
}

type CubeListRequest struct {
	Cubelist []string `json:"list"`
}
