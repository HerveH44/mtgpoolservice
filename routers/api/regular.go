package api

import (
	"mtgpoolservice/logging"
	pack "mtgpoolservice/pool"
	"net/http"

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
		logging.Warn(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ret := make([]*Pack, 0)
	if request.IsDraft {
		for s := 0; s < len(request.Sets); s++ {
			for p := 0; p < request.Players; p++ {
				packs, err := r.packsService.MakeRegularPacks(request.Sets[s : s+1])
				if err != nil {
					logging.Warn(err)
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				ret = append(ret, packs...)
			}
		}
	} else {
		for p := 0; p < request.Players; p++ {
			packs, err := r.packsService.MakeRegularPacks(request.Sets)
			if err != nil {
				logging.Warn(err)
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			ret = append(ret, packs...)
		}
	}

	c.JSON(http.StatusOK, ret)
}

type RegularRequest struct {
	Players int      `json:"players"`
	Sets    []string `json:"sets"`
	IsDraft bool     `json:"isDraft,omitempty"`
}
