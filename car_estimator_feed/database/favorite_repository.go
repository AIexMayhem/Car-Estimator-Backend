package database

import (
    "context"
    "database/sql"
    "fmt"
)

type Favorite struct {
    ID         int64  `db:"id"`
    UserID     string `db:"user_id"`
    ListingID  string `db:"listing_id"`
    CreatedAt  string `db:"created_at"`
}

type FavoriteRepository interface {
    AddToFavorites(ctx context.Context, userID, listingID string) error
}

type favoriteRepoDB struct {
    db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) FavoriteRepository {
    return &favoriteRepoDB{db: db}
}

func (r *favoriteRepoDB) AddToFavorites(ctx context.Context, userID, listingID string) error {
    var exists bool
    err := r.db.QueryRowContext(ctx,
        "SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND listing_id = $2)",
        userID, listingID,
    ).Scan(&exists)
    if err != nil {
        return fmt.Errorf("AddToFavorites check exist failed: %w", err)
    }
    if exists {
        return fmt.Errorf("already in favorites")
    }

    _, err = r.db.ExecContext(ctx,
        "INSERT INTO favorites(user_id, listing_id) VALUES($1, $2)",
        userID, listingID,
    )
    if err != nil {
        return fmt.Errorf("AddToFavorites insert failed: %w", err)
    }
    return nil
}