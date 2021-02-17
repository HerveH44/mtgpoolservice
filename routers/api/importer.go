package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"mtgpoolservice/importer"
)

type ImporterController interface {
	ImportAllSets(context *gin.Context)
	ImportSet(context *gin.Context)
}

type importerController struct {
	importerFacade importer.Facade
}

func NewImporterController(facade importer.Facade) ImporterController {
	return &importerController{facade}
}

func (i *importerController) ImportAllSets(context *gin.Context) {
	err := i.importerFacade.UpdateSets(false)
	if err != nil {
		log.Println(err)
		context.JSON(500, gin.H{"error": "unexpected error"})
	}
	context.JSON(200, gin.H{"success": "ok"})
}

func (i *importerController) ImportSet(context *gin.Context) {
	setCode := context.Param("setCode")
	if setCode == "" {
		context.JSON(500, "unexpected error: Could not fetch the MTGJson monoSet")
	}

	err := i.importerFacade.UpdateSet(setCode)
	if err != nil {
		log.Println(err)
		context.JSON(500, gin.H{"error": "unexpected error"})
	}
	context.JSON(200, gin.H{"success": "ok"})
}

