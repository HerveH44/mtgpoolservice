package db

import (
	"github.com/patrickmn/go-cache"
	"mtgpoolservice/models/entities"
	"time"
)

// Keep the sets for 60 minutes at most...
// Could be more if we have enough memory and updates of MTGJson are not often
var c = cache.New(60*time.Minute, 10*time.Minute)

func GetSet(setCode string) (*entities.Set, error) {
	if cachedSet, found := c.Get(setCode); found {
		return cachedSet.(*entities.Set), nil
	}

	fetchedSet, err := fetchSet(setCode)
	if err != nil {
		return nil, err
	}

	c.SetDefault(setCode, fetchedSet)
	return fetchedSet, nil
}
