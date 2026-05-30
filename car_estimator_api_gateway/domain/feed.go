package domain

import "time"

type CarListing struct {
	ListingId        string    `json:"listing_id"`
	SellerId         string    `json:"seller_id"`
	Description      string    `json:"description"`
	PostedAt         time.Time `json:"posted_at"`
	Status           string    `json:"status"`
	DealType         string    `json:"deal_type"`
	Price            float64   `json:"price"`
	Tags             []string  `json:"tags"`
	CarId            string    `json:"car_id"`
	Mileage          int32     `json:"mileage"`
	OwnersCount      int32     `json:"owners_count"`
	AccidentsCount   int32     `json:"accidents_count"`
	Condition        string    `json:"condition"`
	Color            string    `json:"color"`
	ConfigId         string    `json:"config_id"`
	EngineType       string    `json:"engine_type"`
	EngineVolume     string    `json:"engine_volume"`
	EnginePower      int32     `json:"engine_power"`
	Cylinders        int32     `json:"cylinders"`
	Transmission     string    `json:"transmission"`
	Drivetrain       string    `json:"drivetrain"`
	ModelId          string    `json:"model_id"`
	ModelName        string    `json:"model_name"`
	Make             string    `json:"make"`
	Year             int32     `json:"year"`
	BodyType         string    `json:"body_type"`
	Generation       string    `json:"generation"`
	WeightKg         float64   `json:"weight_kg"`
	SellerName       string    `json:"seller_name"`
	SellerRating     float64   `json:"seller_rating"`
	SellerSalesCount int32     `json:"seller_sales_count"`
	SellerIsBusiness bool      `json:"seller_is_business"`
}

type PageRequest struct {
	PageNumber int32 `json:"page_number"`
	PageSize   int32 `json:"page_size"`
}

type PageResponseMetadata struct {
	TotalItems  int32 `json:"total_items"`
	TotalPages  int32 `json:"total_pages"`
	CurrentPage int32 `json:"current_page"`
}

type ListListingsRequest struct {
	Page   PageRequest `json:"page"`
	SortBy string      `json:"sort_by"`
}

type ListListingsResponse struct {
	Listings     []CarListing       `json:"listings"`
	PageMetadata PageResponseMetadata `json:"page_metadata"`
}

type SearchListingsRequest struct {
	Query  string      `json:"query"`
	Page   PageRequest `json:"page"`
	SortBy string      `json:"sort_by"`
}

type SearchListingsResponse struct {
	Listings     []CarListing       `json:"listings"`
	PageMetadata PageResponseMetadata `json:"page_metadata"`
}

type GetListingResponse struct {
	Listing CarListing `json:"listing"`
}

type CreateListingRequest struct {
	Listing CarListing `json:"listing"`
}

type CreateListingResponse struct {
	Listing CarListing `json:"listing"`
}

type UpdateListingRequest struct {
	Listing CarListing `json:"listing"`
}

type UpdateListingResponse struct {
	Listing CarListing `json:"listing"`
}

type DeleteListingResponse struct {
	Success bool `json:"success"`
}

type AddToFavoritesRequest struct {
	ListingId string `json:"listing_id"`
}

type AddToFavoritesResponse struct {
	Success bool `json:"success"`
}