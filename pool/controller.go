package pool

import (
	"fmt"
	"mtgpoolservice/routers"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type poolController struct {
	packsService Service
}

func NewPoolController(service Service) routers.Controller {
	return &poolController{service}
}

func (r *poolController) Register(router *gin.RouterGroup) {
	router.POST("/regular", r.RegularPacks)
	router.POST("/chaos", r.ChaosPacks)
	router.POST("/cube", r.CubePacks)
	router.POST("/cube/list", r.CubeList)
}

func (r *poolController) RegularPacks(c *gin.Context) {
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

func (r *poolController) ChaosPacks(c *gin.Context) {
	var req ChaosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packs, err := r.packsService.MakeChaosPacks(req.Modern, req.TotalChaos, req.Players*req.Packs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packs)
}

type ChaosRequest struct {
	Players    int  `json:"players"`
	Packs      int  `json:"pool"`
	Modern     bool `json:"modern"`
	TotalChaos bool `json:"totalChaos"`
}

func (r *poolController) CubePacks(c *gin.Context) {
	var req CubeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.PlayerPackSize*req.Players*req.Packs > len(req.Cubelist) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cube list is too small"})
		return
	}

	packs, err := r.packsService.MakeCubePacks(req.Cubelist[:], req.PlayerPackSize, req.Packs*req.Packs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packs)
}

func (r *poolController) CubeList(c *gin.Context) {
	var request CubeListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errs := r.packsService.CheckCubeList(request.Cubelist[:])
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
	Players        int      `json:"players"`
	PlayerPackSize int      `json:"playerPackSize"`
	Packs          int      `json:"packs"`
}

type CubeListRequest struct {
	Cubelist []string `json:"list"`
}
