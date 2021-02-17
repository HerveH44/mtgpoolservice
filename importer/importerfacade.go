package importer

import (
	"log"
	"mtgpoolservice/db"
	"mtgpoolservice/importer/mtgjson"
	"mtgpoolservice/logging"
	"time"
)

type Facade interface {
	UpdateSets(force bool) error
	UpdateSet(code string) error
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
		logging.Info("Update Sets: Found new version", remoteVersion.Data.Version, remoteVersion.Data.Date)
	}

	err := i.mtgjsonService.DownloadSets(i.onSet)
	if err != nil {
		return err
	}

	remoteVersion, err := i.mtgjsonService.DownloadVersion()
	if err != nil {
		return err
	}

	//TODO: we should have all this inside one transaction (if fail, I should rollback everything?)
	version := mtgjson.MapMTGJsonVersionToVersion(remoteVersion)
	return i.versionRepository.SaveVersion(version)
}

func (i *importerFacade) UpdateSet(code string) error {
	mtgJsonSet, err := i.mtgjsonService.DownloadSet(code)
	if err != nil {
		return err
	}
	set := mtgjson.MapMTGJsonSetToEntity(mtgJsonSet, notCubableFunc)
	return i.setRepository.SaveSet(set)
}

func notCubableFunc(string, *mtgjson.Card) bool {
	return false
}

func (i *importerFacade) isNewVersion(version mtgjson.Version) bool {
	getVersion, err := i.versionRepository.GetVersion()
	if err != nil {
		log.Printf("Could not find version in DB because of error %s. Will fallback to version %s as new\n", err, version.Data.Version)
		return true
	}
	return isNewVersion(version, getVersion)
}

func (i *importerFacade) onSet(set *mtgjson.MTGJsonSet, cubable mtgjson.IsCubable) {
	logging.Debug("Mapping set", set.Code)
	mappedSet := mtgjson.MapMTGJsonSetToEntity(set, cubable)
	logging.Debug("Finished mapping set", set.Code)

	if len(mappedSet.Cards) == 0 {
		logging.Debug("Will not save empty set", set.Code)
		return
	}

	if err := i.setRepository.SaveSet(mappedSet); err != nil {
		logging.Debug("Could not save set", set.Code, err)
	} else {
		logging.Debug("Saved set", set.Code)
	}
}

func isNewVersion(remoteVersion mtgjson.Version, savedVersion db.Version) bool {
	remoteVersionDate := time.Time(remoteVersion.Data.Date)
	return savedVersion.Date.Before(remoteVersionDate)
}

func NewImporterFacade(mtgjsonService mtgjson.MTGJsonService, setRepository db.SetRepository, versionRepository db.VersionRepository) Facade {
	return &importerFacade{mtgjsonService, setRepository, versionRepository}
}
