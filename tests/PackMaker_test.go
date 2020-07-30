package tests

import (
	"encoding/json"
	"fmt"
	"mtgpoolservice/common"
	"mtgpoolservice/services"
	"testing"
)

func TestGetPack(t *testing.T) {
	common.Init()
	set, err := common.GetSet("ISD")

	if err != nil {
		t.Error(err)
	}

	pack, err := services.MakePack(&set)

	bytes, err := json.Marshal(pack)

	json := string(bytes)
	fmt.Println(json)

}
