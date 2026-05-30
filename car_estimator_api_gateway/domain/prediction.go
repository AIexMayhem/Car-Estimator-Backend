package domain

type PredictionRequest struct {
	Make     string `json:"make"`
	Model    string `json:"model"`
	Year     int    `json:"year"`
	Hp       int    `json:"hp"`
	Body     string `json:"body"`
	YearSell int    `json:"yearSell"`
	Odometer int    `json:"odometer"`
	Color    string `json:"color"`
}

type PredictionResponse struct {
	Price     int      `json:"price"`
	SellCount int      `json:"sell_count"`
	Urls      []string `json:"urls"`
	GraphImg  string   `json:"graph_img"`
}

type ImageResponse struct {
	Urls []string `json:"urls"`
}
