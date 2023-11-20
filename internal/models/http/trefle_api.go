package httpModels

type SearchSliceResponse struct {
	Data []ItemPlantResponse `json:"data"`
}

type SearchResponse struct {
	Data ItemPlantResponse `json:"data"`
}

type ItemPlantResponse struct {
	ID       int    `json:"id"`
	Slug     string `json:"slug"`
	ScName   string `json:"scientific_name"`
	ImageURL string `json:"image_url"`
}
