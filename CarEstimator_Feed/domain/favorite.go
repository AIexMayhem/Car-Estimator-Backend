package domain

import "time"

type Favorite struct {
    UserID    string
    ListingID string
    CreatedAt time.Time
}