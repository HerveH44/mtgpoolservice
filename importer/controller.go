package importer

import (
	"mtgpoolservice/routers"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type importerController struct {
	importerFacade Facade
}

func (i *importerController) Register(router *gin.RouterGroup) {
	router.POST("", i.importAllSets)
}

func NewImporterController(facade Facade) routers.Controller {
	return &importerController{facade}
}

func (i *importerController) importAllSets(context *gin.Context) {
	context.JSON(200, gin.H{"success": "ok"})
	go i.doImportAllSets()
}

func (i *importerController) doImportAllSets() {
	err := i.importerFacade.UpdateSets(false)
	if err != nil {
		log.Error("Could not import sets ", err)
	}
}
