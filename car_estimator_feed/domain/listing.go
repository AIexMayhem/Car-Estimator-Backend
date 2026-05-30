package domain

import "time"

type Listing struct {
    ListingID        string
    SellerID         string
    Description      string
    PostedAt         time.Time
    Status           string
    DealType         string
    Price            float64
    CarID            string
    Mileage          int
    OwnersCount      int
    AccidentsCount   int
    Condition        string
    Color            string
    ConfigID         string
    EngineType       string
    EngineVolume     string
    EnginePower      int
    Cylinders        int
    Transmission     string
    Drivetrain       string
    ModelID          string
    ModelName        string
    Make             string
    Year             int
    BodyType         string
    Generation       string
    WeightKg         float64
    SellerName       string
    SellerRating     float64
    SellerSalesCount int
    SellerIsBusiness bool
}
