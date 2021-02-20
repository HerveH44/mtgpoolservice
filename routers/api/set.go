package api

import (
	database "mtgpoolservice/db"
	"net/http"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

type SetController interface {
	GetAvailableSets(context *gin.Context)
	GetLatestSet(context *gin.Context)
	GetInfos(context *gin.Context)
	GetVersion(context *gin.Context)
}

type setController struct {
	setRepo     database.SetRepository
	versionRepo database.VersionRepository
	cache       *cache.Cache
}

var getInfosKey = "__infos__"

func NewSetController(setRepository database.SetRepository, versionRepository database.VersionRepository) SetController {
	setControllerCache := cache.New(6*time.Hour, 1*time.Hour)
	return &setController{setRepository, versionRepository, setControllerCache}
}

func (s *setController) GetInfos(context *gin.Context) {
	if cachedInfos, found := s.cache.Get(getInfosKey); found {
		context.JSON(200, cachedInfos)
	}

	version, err := s.versionRepo.GetVersion()
	if err != nil {
		log.Error(err)
		context.JSON(500, gin.H{"error": "unexpected error while fetching infos about version"})
		return
	}

	setMap, err := s.getAvailableSets()
	if err != nil {
		log.Error(err)
		context.JSON(500, gin.H{"error": "unexpected error while fetching available sets"})
		return
	}

	latestSet, err := s.getLatestSet()
	if err != nil {
		log.Error(err)
		context.JSON(500, gin.H{"error": "unexpected error while fetching latest set"})
		return
	}

	response := Infos{
		AvailableSetsMap:  setMap,
		LatestSetResponse: latestSet,
		MTGJsonVersion: VersionResponse{
			Date:    version.Date.Format("2006-01-02"),
			Version: version.SemanticVersion,
		},
		BoosterRulesVersion: version.SemanticVersion,
	}
	s.cache.SetDefault(getInfosKey, response)

	context.JSON(200, response)
}

func (s *setController) GetVersion(context *gin.Context) {
	log.Info("Getting version")

	version, err := s.versionRepo.GetVersion()
	if err != nil {
		log.Error(err)
		context.JSON(500, gin.H{"error": "unexpected error"})
	}

	response := VersionResponse{
		Date:    version.Date.Format("2006-01-02"),
		Version: version.SemanticVersion,
	}

	context.JSON(200, response)
}

func (s *setController) GetAvailableSets(context *gin.Context) {
	setMap, err := s.getAvailableSets()
	if err != nil {
		context.JSON(500, gin.H{"error": "unexpected error"})
		return
	}

	context.JSON(http.StatusOK, setMap)
}

func (s *setController) GetLatestSet(context *gin.Context) {
	latestSet, err := s.getLatestSet()
	if err != nil {
		context.JSON(500, gin.H{"error": "unexpected error"})
		return
	}

	context.JSON(http.StatusOK, latestSet)
}

func (s *setController) getAvailableSets() (AvailableSetsMap, error) {
	sets, err := s.setRepo.FindAllSets()
	if err != nil {
		return nil, err
	}
	return buildAvailableSetsMap(sets), nil
}

func (s *setController) getLatestSet() (LatestSetResponse, error) {
	latestSet, err := s.setRepo.FindLatestSet()
	if err != nil {
		return LatestSetResponse{}, err
	}

	return LatestSetResponse{
		SetResponse: SetResponse{
			Code: latestSet.Code,
			Name: latestSet.Name,
		},
		Type: latestSet.Type,
	}, nil
}

func buildAvailableSetsMap(sets []*database.Set) AvailableSetsMap {
	setMap := setupSetsMap()

	for _, set := range sets {
		setType := set.Type
		i := sort.SearchStrings(playableSetTypes, setType)
		if i >= len(playableSetTypes) || playableSetTypes[i] != setType {
			continue
		}
		setMap[setType] = append(setMap[setType], SetResponse{
			Code: set.Code,
			Name: set.Name,
		})
	}
	return setMap
}

func setupSetsMap() (setMap AvailableSetsMap) {
	setMap = make(AvailableSetsMap)
	for _, t := range playableSetTypes {
		setMap[t] = make([]SetResponse, 0)
	}

	setMap["random"] = make([]SetResponse, 0)
	setMap["random"] = append(setMap["random"], SetResponse{
		Code: "RNG",
		Name: "Random Set",
	})
	sort.Strings(playableSetTypes)

	return setMap
}

var playableSetTypes = []string{"core", "expansion", "draft_innovation", "funny", "starter", "masters"}

type LatestSetResponse struct {
	SetResponse
	Type string `json:"type"`
}

type VersionResponse struct {
	Date    string `json:"date"`
	Version string `json:"version"`
}

type AvailableSetsMap map[string][]SetResponse

type Infos struct {
	AvailableSetsMap    AvailableSetsMap  `json:"availableSets"`
	LatestSetResponse   LatestSetResponse `json:"latestSet"`
	MTGJsonVersion      VersionResponse   `json:"mtgJsonVersion"`
	BoosterRulesVersion string            `json:"boosterRulesVersion"`
}

type SetResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
