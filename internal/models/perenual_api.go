package models

type PerenualsPlantResponse struct {
	Data []PerenualPlant
}

type PerenualPlantDetail struct {
	ID                       uint64            `json:"id"`
	CommonName               string            `json:"common_name"`
	ScientificName           []string          `json:"scientific_name"`
	OtherName                []string          `json:"other_name"`
	Family                   string            `json:"family"`
	Origin                   string            `json:"origin"`
	Type                     string            `json:"type"`
	Dimensions               Dimensions        `json:"dimensions"`
	Cycle                    string            `json:"cycle"`
	Watering                 string            `json:"watering"`
	WateringPeriod           string            `json:"watering_period"`
	DepthWaterReq            UnitValue         `json:"depth_water_requirement"`
	VolumeWaterReq           UnitValue         `json:"volume_water_requirement"`
	WateringGeneralBenchmark UnitValue         `json:"watering_general_benchmark"`
	PlantAnatomy             PlantAnatomy      `json:"plant_anatomy"`
	Sunlight                 []string          `json:"sunlight"`
	PruningCount             PuringCount       `json:"puring_count"`
	PruningMonth             []string          `json:"pruning_month"`
	Seeds                    uint64            `json:"seeds"`
	Attracts                 []string          `json:"attracts"`
	Propagation              []string          `json:"propagation"`
	Hardiness                Hardiness         `json:"hardiness"`
	HardinessLocation        HardinessLocation `json:"hardiness_location"`
	Flowers                  bool              `json:"flowers"`
	FloweringSeason          string            `json:"flowering_season"`
	Clolor                   string            `json:"color"`

	Cones            bool     `json:"cones"`
	Fruits           bool     `json:"fruits"`
	EdibleFruit      bool     `json:"edible_fruit"`
	Leaf             bool     `json:"leaf"`
	LeafColor        []string `json:"leaf_color"`
	GrowthRate       string   `json:"growth_rate"`
	Maintance        string   `json:"maintance"`
	PoisonousToHumas bool     `json:"poisonous_to_humans"`
	PoisonousToPets  bool     `json:"poisonous_to_pets"`

	ImageURL ImageURL `json:"default_image"`
}

type PerenualPlant struct {
	ID             uint64   `json:"id"`
	CommonName     string   `json:"common_name"`
	ScientificName []string `json:"scientific_name"`
	OtherName      []string `json:"other_name"`
	Cycle          string   `json:"cycle"`
	Watering       string   `json:"watering"`
	Sunlight       []string `json:"sunlight"`
	Premium        string   `json:"sunlight"`
	ImageURL       ImageURL `json:"default_image"`
}

type ImageURL struct {
	URL string `json:"original_url"`
}

type PuringCount struct {
	Amount   uint64 `json:"amount"`
	Interval string `json:"string"`
}

type UnitValue struct {
	Unit  string `json:"unit"`
	Value string `json:"value"`
}

type Dimensions struct {
	Type     string  `json:"type"`
	MinValue float64 `json:"min_value"`
	MaxValue float64 `json:"max_value"`
	Unit     string  `json:"unit"`
}

type PlantAnatomy struct {
	Bark   string `json:"bark"`
	Leaves string `json:"leaves"`
}

type Hardiness struct {
	Min uint64 `json:"min"`
	Max uint64 `json:"max"`
}

type HardinessLocation struct {
	FullURL    string `json:"full_url"`
	FullIframe string `json:"full_iframe"`
}
