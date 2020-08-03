package api

import (
	"github.com/gin-gonic/gin"
	"mtgpoolservice/db"
	"mtgpoolservice/logging"
	"mtgpoolservice/models"
)

func Version(context *gin.Context) {
	logging.Info("Getting version")

	version, err := db.FetchLastVersion()
	if err != nil {
		logging.Error(err)
		context.JSON(500, gin.H{"error": "unexpected error"})
	}

	response := models.VersionResponse{
		Date:    version.Date.Format("2006-01-02"),
		Version: version.SemanticVersion,
	}

	context.JSON(200, response)
}
