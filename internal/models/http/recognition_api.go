package httpModels

type Project string

const (
	AllProject Project = "all"
)

type Result struct {
	Score   float64 `json:"score"`
	Species struct {
		Name string `json:"scientificNameWithoutAuthor"`
	} `json:"species"`
}

type RecResponse struct {
	Results []Result `json:"results"`
}
