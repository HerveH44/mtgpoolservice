package db

import (
	"mtgpoolservice/routers"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

const getInfosKey = "__infos__"

var playableSetTypes = []string{"core", "expansion", "draft_innovation", "funny", "starter", "masters"}

type dbController struct {
	setRepo     SetRepository
	versionRepo VersionRepository
	cache       *cache.Cache
}

func (s *dbController) Register(router *gin.RouterGroup) {
	router.GET("", s.getInfos)
}

func NewDBController(setRepository SetRepository, versionRepository VersionRepository) routers.Controller {
	setControllerCache := cache.New(6*time.Hour, 1*time.Hour)
	return &dbController{setRepository, versionRepository, setControllerCache}
}

func (s *dbController) getInfos(context *gin.Context) {
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

func (s *dbController) getAvailableSets() (AvailableSetsMap, error) {
	sets, err := s.setRepo.FindAllSets()
	if err != nil {
		return nil, err
	}
	return buildAvailableSetsMap(sets), nil
}

func (s *dbController) getLatestSet() (LatestSetResponse, error) {
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

func buildAvailableSetsMap(sets []*Set) AvailableSetsMap {
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
