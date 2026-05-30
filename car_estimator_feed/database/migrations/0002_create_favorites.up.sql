CREATE TABLE IF NOT EXISTS favorites (
    id          SERIAL PRIMARY KEY,
    user_id     VARCHAR(36) NOT NULL,
    listing_id  VARCHAR(36) NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_listing
        FOREIGN KEY (listing_id)
        REFERENCES listings(listing_id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_favorites_user
    ON favorites(user_id);