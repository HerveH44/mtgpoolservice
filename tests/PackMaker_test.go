package tests

import (
	"encoding/json"
	"fmt"
	"mtgpoolservice/db"
	"mtgpoolservice/services"
	"testing"
)

func TestGetPack(t *testing.T) {
	db.Init()
	set, err := db.GetSet("ISD")

	if err != nil {
		t.Error(err)
	}

	pack, err := services.MakePack(&set)

	bytes, err := json.Marshal(pack)

	json := string(bytes)
	fmt.Println(json)

}
