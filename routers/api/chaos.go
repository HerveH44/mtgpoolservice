package api

import (
	pack "mtgpoolservice/pool"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChaosController interface {
	ChaosPacks(c *gin.Context)
}

type chaosController struct {
	packService pack.Service
}

func NewChaosController(service pack.Service) ChaosController {
	return &chaosController{service}
}

func (cc *chaosController) ChaosPacks(c *gin.Context) {
	var req ChaosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packs, err := cc.packService.MakeChaosPacks(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packs)
}

type ChaosRequest struct {
	Players    uint `json:"players"`
	Packs      uint `json:"pool"`
	Modern     bool `json:"modern"`
	TotalChaos bool `json:"totalChaos"`
}
