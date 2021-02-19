package api

import (
	pack "mtgpoolservice/pool"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type RegularController interface {
	RegularPacks(c *gin.Context)
}

type regularController struct {
	packsService pack.Service
}

func NewRegularController(service pack.Service) RegularController {
	return &regularController{service}
}

func (r *regularController) RegularPacks(c *gin.Context) {
	var request RegularRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packs, err := r.packsService.MakeRegularPacks(request.Sets, request.Players)
	if err != nil {
		log.Warn(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packs)
}

type RegularRequest struct {
	Players int      `json:"players"`
	Sets    []string `json:"sets"`
}
