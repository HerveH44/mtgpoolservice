package models

import (
	"github.com/google/uuid"
	"math/rand"
	"mtgpoolservice/models/entities"
	"time"
)

type SetResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type VersionResponse struct {
	Date    string `json:"date"`
	Version string `json:"version"`
}

type RegularRequest struct {
	Players int      `json:"players"`
	Sets    []string `json:"sets"`
}

type CardResponse struct {
	entities.Card

	Id   string `json:"cardId"`
	Foil bool   `json:"foil"`
}

type CardPool []CardResponse

// Shuffle shuffles the array parameter in place
func (c *CardPool) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(*c), func(i, j int) { (*c)[i], (*c)[j] = (*c)[j], (*c)[i] })
}

func (c *CardPool) Add(card *entities.Card, isFoil bool) {
	cardResponse := CardResponse{
		Card: *card,
		Id:   uuid.New().String(),
		Foil: isFoil,
	}
	*c = append(*c, cardResponse)
}

func (c *CardPool) AddCards(cards *[]entities.Card) {
	for _, card := range *cards {
		c.Add(&card, false)
	}
}

type ChaosRequest struct {
	Players    uint `json:"players"`
	Packs      uint `json:"packs"`
	Modern     bool `json:"modern"`
	TotalChaos bool `json:"totalChaos"`
}

type CubeRequest struct {
	Cubelist       []string `json:"list"`
	Players        uint     `json:"players"`
	PlayerPackSize uint     `json:"playerPackSize"`
	Packs          uint     `json:"packs"`
}

type CubeListRequest struct {
	Cubelist []string `json:"list"`
}

type AvailableSetsMap map[string][]SetResponse
