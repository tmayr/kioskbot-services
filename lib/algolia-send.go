package lib

import (
	KioskTypes "kioskbot-services/types"
	"log"
	"os"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

func SendProductsToAlgolia(items []KioskTypes.Item) {
	client := algoliasearch.NewClient(os.Getenv("ALGOLIA_ID"), os.Getenv("ALGOLIA_KEY"))
	index := client.InitIndex(os.Getenv("ALGOLIA_PRODUCTS_INDEX"))

	var algoliaObjects []algoliasearch.Object
	var algoliaSynonyms []algoliasearch.Synonym

	for _, v := range items {
		slug := v.Slug
		item := v.Item
		ID := v.ID
		synonyms := v.Synonyms

		algoliaObjects = append(algoliaObjects, algoliasearch.Object{
			"_id":      ID,
			"item":     item,
			"slug":     slug,
			"objectID": slug,
		})

		if synonyms != nil {
			for _, v := range synonyms {
				algoliaSynonyms = append(
					algoliaSynonyms,
					algoliasearch.NewOneWaySynonym(slug, v, []string{
						item,
					}),
				)
			}
		}
	}

	_, err := index.AddObjects(algoliaObjects)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = index.BatchSynonyms(algoliaSynonyms, true, true)
	if err != nil {
		log.Fatal(err)
		return
	}
}
