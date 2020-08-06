package db

import (
	"github.com/patrickmn/go-cache"
	"mtgpoolservice/models/entities"
	"time"
)

var getSetsKey = "__available_sets__"

// Keep the sets for 60 minutes at most...
// Could be more if we have enough memory and updates of MTGJson are not often
var setCache = cache.New(60*time.Minute, 10*time.Minute)
var cardCache = cache.New(60*time.Minute, 10*time.Minute)
var unknownCardCache = cache.New(60*time.Minute, 10*time.Minute)

func GetSet(setCode string) (*entities.Set, error) {
	if cachedSet, found := setCache.Get(setCode); found {
		return cachedSet.(*entities.Set), nil
	}

	fetchedSet, err := fetchSet(setCode)
	if err != nil {
		return nil, err
	}

	setCache.SetDefault(setCode, fetchedSet)
	return fetchedSet, nil
}

func GetSets() (*[]entities.Set, error) {
	if cachedSets, found := setCache.Get(getSetsKey); found {
		return cachedSets.(*[]entities.Set), nil
	}

	fetchedSets, err := getSets()
	if err != nil {
		return nil, err
	}

	setCache.SetDefault(getSetsKey, fetchedSets)
	return fetchedSets, nil
}

func GetCardWithName(name string) (card entities.Card, err error) {
	if foundCard, found := cardCache.Get(name); found {
		return foundCard.(entities.Card), nil
	}

	if foundCard, found := unknownCardCache.Get(name); found {
		return card, foundCard.(error)
	}

	card, err = getCardWithName(name)
	if err != nil {
		unknownCardCache.SetDefault(name, err)
		return card, err
	}

	cardCache.SetDefault(name, card)
	return card, nil
}
