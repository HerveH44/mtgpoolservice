package models

import "mtgpoolservice/models/entities"

type RegularRequest struct {
	Players int      `json:"players"`
	Sets    []string `json:"sets"`
}

type RegularDraftResponse [][]Pool

type CardResponse struct {
	*entities.Card

	Id   string `json:"id"`
	Foil bool   `json:"foil"`
}

type Pool []CardResponse

type ChaosRequest struct {
	Players    uint `json:"players`
	Packs      uint `json:"packs"`
	Modern     bool `json:"modern"`
	TotalChaos bool `json:"total_chaos"`
}

type CubeDraftRequest struct {
	Cubelist       []string `json:"list"`
	Players        uint     `json:"players`
	PlayerPackSize uint     `json:"player_pack_size"`
	Packs          uint     `json:"packs"`
}

type CubeSealedRequest struct {
	Cubelist       []string `json:"list"`
	Players        uint     `json:"players`
	PlayerPoolSize uint     `json:"player_pool_size"`
}
