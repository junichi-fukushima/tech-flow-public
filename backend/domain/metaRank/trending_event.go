package metaRank

type Trending struct {
	Took  int          `json:"took"`
	Items []ItemDetail `json:"items"`
}

type ItemDetail struct {
	Item  string  `json:"item"`
	Score float64 `json:"score"`
}
