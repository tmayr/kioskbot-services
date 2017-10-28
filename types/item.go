package types

type Item struct {
	Item      string   `json:"item"`
	Slug      string   `json:"slug"`
	Synonyms  []string `json:"synonyms"`
	CreatedAt string   `json:"createdAt"`
}
