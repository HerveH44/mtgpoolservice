package models

import "mtgpoolservice/models/entities"

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

type RegularDraftResponse [][]Pool

type CardResponse struct {
	entities.Card

	Id   string `json:"cardId"`
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
	PlayerPackSize uint     `json:"playerPackSize"`
	Packs          uint     `json:"packs"`
}

type CubeListRequest struct {
	Cubelist []string `json:"list"`
}

type CubeSealedRequest struct {
	Cubelist       []string `json:"list"`
	Players        uint     `json:"players`
	PlayerPoolSize uint     `json:"player_pool_size"`
}
