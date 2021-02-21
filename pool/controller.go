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
	router.POST("/regular/draft", r.regularDraft)
	router.POST("/regular/sealed", r.regularSealed)
	router.POST("/chaos/draft", r.chaosDraft)
	router.POST("/chaos/sealed", r.chaosSealed)
	router.POST("/cube", r.cubePacks)
	router.POST("/cube/list", r.cubeList)
}

func (r *poolController) regularDraft(c *gin.Context) {
	var request RegularRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packs, err := r.packsService.regularDraft(request.Sets, request.Players)
	if err != nil {
		log.Warn(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packs)
}

func (r *poolController) regularSealed(c *gin.Context) {
	var request RegularRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Warn(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packs, err := r.packsService.regularSealed(request.Sets, request.Players)
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

func (r *poolController) chaosDraft(c *gin.Context) {
	var req ChaosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packs, err := r.packsService.chaosDraft(req.Modern, req.TotalChaos, req.Players*req.Packs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packs)
}

func (r *poolController) chaosSealed(c *gin.Context) {
	var req ChaosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	packs, err := r.packsService.chaosSealed(req.Modern, req.TotalChaos, req.Packs, req.Players)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packs)
}

type ChaosRequest struct {
	Players    int  `json:"players"`
	Packs      int  `json:"packs"`
	Modern     bool `json:"modern"`
	TotalChaos bool `json:"totalChaos"`
}

func (r *poolController) cubePacks(c *gin.Context) {
	var req CubeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.PlayerPackSize*req.Players*req.Packs > len(req.Cubelist) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cube list is too small"})
		return
	}

	packs, err := r.packsService.cubePacks(req.Cubelist[:], req.PlayerPackSize, req.Packs*req.Packs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packs)
}

func (r *poolController) cubeList(c *gin.Context) {
	var request CubeListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errs := r.packsService.checkCubeList(request.Cubelist[:])
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
