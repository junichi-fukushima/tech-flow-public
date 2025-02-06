package metaRank

type Recommend struct {
	Took  int                   `json:"took"`
	Items []RecommendItemDetail `json:"items"`
}

type RecommendItemDetail struct {
	Item     string    `json:"item"`
	Score    float64   `json:"score"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}
