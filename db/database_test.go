package db

import (
	"testing"
)

func TestGetSet(t *testing.T) {
	Init()
	set, err := GetSet("ISD")

	if err != nil {
		t.Error(err)
	}

	if set.Code != "ISD" {
		t.Error("expected ISD set")
	}
}
