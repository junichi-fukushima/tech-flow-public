package dto

type TrendingMetadataEventRequest struct {
	Count int `json:"count"`
}

func (c *TrendingMetadataEventRequest) FromTrendingMetadataEvent(limit int) *TrendingMetadataEventRequest {
	return &TrendingMetadataEventRequest{
		Count: limit,
	}
}
