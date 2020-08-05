package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"mtgpoolservice/services"
	"net/http"
)

func AvailableSets(context *gin.Context) {
	setMap, err := services.GetAvailableSets()
	if err != nil {
		context.JSON(500, gin.H{"error": "unexpected error"})
		return
	}

	context.JSON(http.StatusOK, setMap)
}

func ImportAllSets(context *gin.Context) {
	err := services.UpdateSets()
	if err != nil {
		log.Println(err)
		context.JSON(500, gin.H{"error": "unexpected error"})
	}
	context.JSON(200, gin.H{"success": "ok"})
}

func ImportSet(context *gin.Context) {
	setCode := context.Param("setCode")
	if setCode == "" {
		context.JSON(500, "unexpected error: Could not fetch the MTGJson monoSet")
	}

	err := services.UpdateSet(setCode)
	if err != nil {
		log.Println(err)
		context.JSON(500, gin.H{"error": "unexpected error"})
	}
	context.JSON(200, gin.H{"success": "ok"})
}
