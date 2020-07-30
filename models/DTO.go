package models

type RegularDraftRequest struct {
	Players int      `json:"players"`
	Sets    []string `json:"sets"`
}

type RegularDraftResponse [][]Pack

type CardResponse struct {
	Card

	Id   string `json:"id"`
	Foil bool   `json:"foil"`
}

type Pack []CardResponse
