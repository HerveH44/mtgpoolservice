package importer

import (
	"mtgpoolservice/db"
	"mtgpoolservice/importer/mtgjson"
	"mtgpoolservice/utils"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
)

type Facade interface {
	UpdateSets(force bool) error
}

type importerFacade struct {
	mtgjsonService    mtgjson.MTGJsonService
	setRepository     db.SetRepository
	versionRepository db.VersionRepository
}

func (i *importerFacade) UpdateSets(forceUpdate bool) error {
	if !forceUpdate {
		remoteVersion, err := i.mtgjsonService.DownloadVersion()
		if err != nil {
			return err
		}
		if !i.isNewVersion(remoteVersion) {
			log.Println("Remote MTGJson version is same as db. Not updating")
			return nil
		}
		log.Info("Update Sets: Found new version ", remoteVersion.Data.Version, remoteVersion.Data.Date)
	}

	sets, err := i.mtgjsonService.DownloadSets()
	if err != nil {
		return err
	}

	err = i.importSets(sets.Sets)
	if err != nil {
		return err
	}

	version := mtgjson.MapMTGJsonVersionToVersion(sets.Meta)
	return i.versionRepository.SaveVersion(version)
}

func (i *importerFacade) isNewVersion(version mtgjson.Version) bool {
	getVersion, err := i.versionRepository.GetVersion()
	if err != nil {
		log.Printf("Could not find version in DB because of error %s. Will fallback to version %s as new\n", err, version.Data.Version)
		return true
	}
	return isNewVersion(version, getVersion)
}

func (i *importerFacade) importSets(sets map[string]mtgjson.MTGJsonSet) error {
	ordereredSetsCode := orderSetsByDate(&sets)
	isCubable := buildIsCubableFunc(ordereredSetsCode)

	var mappedSets = make([]db.Set, 0)

	for _, set := range sets {
		log.Debug("Mapping set ", set.Code)
		mappedSet := mtgjson.MapMTGJsonSetToEntity(&set, isCubable)
		log.Debug("Finished mapping set ", set.Code)

		if len(mappedSet.Cards) == 0 {
			log.Debug("Not saving empty set ", set.Code)
		} else {
			mappedSets = append(mappedSets, *mappedSet)
		}
	}

	return i.setRepository.SaveSets(&mappedSets)
}

func buildIsCubableFunc(orderedSetCodes map[string]int) func(string, *mtgjson.Card) bool {
	return func(setCode string, card *mtgjson.Card) bool {
		if len(card.Variations) > 0 && (card.BorderColor == "borderless" || card.IsStarter || utils.Include(card.FrameEffects, "extendedart")) {
			return false
		}
		printings := card.Printings
		if len(printings) < 2 {
			return true
		}

		sort.SliceStable(printings, func(i, j int) bool {
			setCode1 := printings[i]
			setCode2 := printings[j]
			set1Index, ok := orderedSetCodes[setCode1]
			if !ok {
				return false
			}
			set2Index, ok := orderedSetCodes[setCode2]
			if !ok {
				return true
			}
			return set1Index < set2Index
		})

		if setCode == printings[0] {
			return true
		}

		return false
	}
}

var notCubableSetTypes = []string{"box", "duel_deck", "masterpiece", "memorabilia", "promo", "spellbook"}

func orderSetsByDate(m *map[string]mtgjson.MTGJsonSet) map[string]int {
	r := make([]string, 0)
	for setCode, set := range *m {
		if !utils.Include(notCubableSetTypes, set.Type) {
			r = append(r, setCode)
		}
	}
	myMap := *m
	sort.SliceStable(r, func(i, j int) bool {
		releaseDate1 := myMap[r[i]].ReleaseDate
		releaseDate2 := myMap[r[j]].ReleaseDate
		return releaseDate1.Before(&releaseDate2)
	})

	ret := make(map[string]int)
	for i, code := range r {
		ret[code] = i
	}
	return ret
}

func isNewVersion(remoteVersion mtgjson.Version, savedVersion db.Version) bool {
	remoteVersionDate := time.Time(remoteVersion.Data.Date)
	return savedVersion.Date.Before(remoteVersionDate)
}

func NewImporterFacade(mtgjsonService mtgjson.MTGJsonService, setRepository db.SetRepository, versionRepository db.VersionRepository) Facade {
	return &importerFacade{mtgjsonService, setRepository, versionRepository}
}
