package models

type Meta struct {
	Date    string `json:"date"`
	Version string `json:"version"`
}

type AllPrintings struct {
	Data map[string]Set `json:"data"`
	Meta Meta           `json:"meta"`
}
