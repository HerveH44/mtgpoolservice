package api

import (
	"mtgpoolservice/importer"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ImporterController interface {
	ImportAllSets(context *gin.Context)
}

type importerController struct {
	importerFacade importer.Facade
}

func NewImporterController(facade importer.Facade) ImporterController {
	return &importerController{facade}
}

func (i *importerController) ImportAllSets(context *gin.Context) {
	context.JSON(200, gin.H{"success": "ok"})
	go i.importAllSets()
}

func (i *importerController) importAllSets() {
	err := i.importerFacade.UpdateSets(false)
	if err != nil {
		log.Error("Could not import sets ", err)
	}
}
