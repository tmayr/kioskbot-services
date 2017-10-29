package types

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Item struct {
	ID          bson.ObjectId `bson:"_id"`
	Item        string        `bson:"item"`
	Slug        string        `bson:"slug"`
	StockActual int           `bson:"stockActual"`
	Synonyms    []string      `bson:"synonyms"`
	CreatedAt   time.Time     `bson:"createdAt"`
}
