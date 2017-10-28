package lib

import (
	KioskTypes "kioskbot-services/types"
	"os"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
)

func SendProductsToAlgolia(items []KioskTypes.Item) {
	client := algoliasearch.NewClient(os.Getenv("ALGOLIA_ID"), os.Getenv("ALGOLIA_KEY"))
	index := client.InitIndex("products_from_DAT_GO")

	var algoliaObjects []algoliasearch.Object
	var algoliaSynonyms []algoliasearch.Synonym

	for _, v := range items {
		slug := v.Slug
		item := v.Item
		synonyms := v.Synonyms

		algoliaObjects = append(algoliaObjects, algoliasearch.Object{
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

	index.AddObjects(algoliaObjects)
	index.BatchSynonyms(algoliaSynonyms, true, true)
}
