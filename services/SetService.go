package services

import (
	database "mtgpoolservice/db"
	"mtgpoolservice/models"
	"mtgpoolservice/models/entities"
	"sort"
)

func GetAvailableSets() (models.AvailableSetsMap, error) {
	sets, err := database.GetSets()
	if err != nil {
		return nil, err
	}
	return buildAvailableSetsMap(sets), nil
}

func GetLatestSet() (models.LatestSetResponse, error) {
	latestSet, err := database.GetLatestSet()
	if err != nil {
		return models.LatestSetResponse{}, err
	}

	return models.LatestSetResponse{
		SetResponse: models.SetResponse{
			Code: latestSet.Code,
			Name: latestSet.Name,
		},
		Type: latestSet.Type,
	}, nil
}

func buildAvailableSetsMap(sets *[]entities.Set) models.AvailableSetsMap {
	setMap := setupSetsMap()

	for _, set := range *sets {
		setType := set.Type
		i := sort.SearchStrings(database.PlayableSetTypes, setType)
		if i >= len(database.PlayableSetTypes) || database.PlayableSetTypes[i] != setType {
			continue
		}
		setMap[setType] = append(setMap[setType], models.SetResponse{
			Code: set.Code,
			Name: set.Name,
		})
	}
	return setMap
}

func setupSetsMap() (setMap models.AvailableSetsMap) {
	setMap = make(models.AvailableSetsMap)
	for _, t := range database.PlayableSetTypes {
		setMap[t] = make([]models.SetResponse, 0)
	}

	setMap["random"] = make([]models.SetResponse, 0)
	setMap["random"] = append(setMap["random"], models.SetResponse{
		Code: "RNG",
		Name: "Random Set",
	})
	sort.Strings(database.PlayableSetTypes)

	return setMap
}

func getChaosSets(modernOnly bool) (*[]entities.Set, error) {
	sets, err := database.GetSets()
	if err != nil {
		return nil, err
	}

	modernOnlySets := make([]entities.Set, 0)
	for _, set := range *sets {
		if !set.IsExpansionOrCore() {
			continue
		}
		if modernOnly && !set.IsModern() {
			continue
		}
		modernOnlySets = append(modernOnlySets, set)
	}
	return &modernOnlySets, nil
}
