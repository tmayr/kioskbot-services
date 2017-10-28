package lib

import (
	KioskTypes "kioskbot-services/types"
	"os"

	mgo "gopkg.in/mgo.v2"
)

func FetchProductsFromMongo() []KioskTypes.Item {
	session, err := mgo.Dial(os.Getenv("MONGO_URI"))
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var results []KioskTypes.Item
	c := session.DB(os.Getenv("MONGO_DB")).C("products")
	c.Find(nil).All(&results)
	return results
}
