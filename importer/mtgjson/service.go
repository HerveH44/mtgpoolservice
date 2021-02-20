package mtgjson

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type MTGJsonService interface {
	DownloadVersion() (Version, error)
	DownloadSets() (*MTGJsonSets, error)
	DownloadSet(setCode string) (*MTGJsonSet, error)
}

type IsCubable func(string, *Card) bool

type mtgJsonService struct {
	endpoint string
}

func NewMTGJsonService(endpoint string) MTGJsonService {
	return &mtgJsonService{endpoint: endpoint}
}

func (m *mtgJsonService) DownloadVersion() (version Version, err error) {
	log.Info("Check MTGJSON remote version")

	resp, err := http.Get(m.endpoint + "Meta.json")
	if err != nil {
		log.Info("Could not fetch the MTGJson version")
		return
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&version); err != nil {
		log.Error("error while unmarshalling version ", err)
	}

	return
}

func (m *mtgJsonService) DownloadSet(setCode string) (*MTGJsonSet, error) {
	resp, err := http.Get(fmt.Sprint(m.endpoint, strings.ToUpper(setCode), ".json"))
	if err != nil {
		return nil, fmt.Errorf("RefreshSet: Could not fetch the MTGJson set with code %s", setCode)
	}
	defer resp.Body.Close()

	monoSet := new(MonoSet)
	if err := json.NewDecoder(resp.Body).Decode(monoSet); err != nil {
		return nil, fmt.Errorf("RefreshSet: Error while unmarshalling monoSet %w", err)
	}

	return &monoSet.Data, nil
}

func (m *mtgJsonService) DownloadSets() (*MTGJsonSets, error) {
	log.Info("Refreshing sets")
	resp, err := http.Get(m.endpoint + "MTGJsonSets.json")
	if err != nil {
		log.Error("Could not fetch the MTGJson sets")
		return nil, err
	}
	defer resp.Body.Close()

	sets := new(MTGJsonSets)
	if err := json.NewDecoder(resp.Body).Decode(sets); err != nil {
		log.Error("error while unmarshalling sets", err)
		return nil, err
	}
	var setsNumber = len(sets.Sets)
	log.Debug(setsNumber, " sets found")
	return sets, nil
}
