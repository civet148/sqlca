package parser

import (
	"go.mongodb.org/mongo-driver/bson"
)

type MgoResult struct {
	Type       SqlType `json:"type"`
	Match      bson.M  `json:"match"`
	Sort       bson.M  `json:"sort"`
	Set        bson.M  `json:"set"`
	Group      bson.M  `json:"group"`
	Projection bson.M  `json:"projection"`
}

func NewMgoResult(typ SqlType) *MgoResult {
	if !typ.IsValid() {
		panic("sql type not valid")
	}
	return &MgoResult{
		Type:       typ,
		Match:      make(bson.M, 0),
		Sort:       make(bson.M, 0),
		Set:        make(bson.M, 0),
		Group:      make(bson.M, 0),
		Projection: make(bson.M, 0),
	}
}
