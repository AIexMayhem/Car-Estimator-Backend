package database

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "time"
)

type Listing struct {
    ListingID        string         `db:"listing_id"`
    SellerID         string         `db:"seller_id"`
    Description      sql.NullString `db:"description"`
    PostedAt         time.Time      `db:"posted_at"`
    Status           sql.NullString `db:"status"`
    DealType         sql.NullString `db:"deal_type"`
    Price            sql.NullFloat64`db:"price"`
    CarID            sql.NullString `db:"car_id"`
    Mileage          sql.NullInt64  `db:"mileage"`
    OwnersCount      sql.NullInt64  `db:"owners_count"`
    AccidentsCount   sql.NullInt64  `db:"accidents_count"`
    Condition        sql.NullString `db:"condition"`
    Color            sql.NullString `db:"color"`
    ConfigID         sql.NullString `db:"config_id"`
    EngineType       sql.NullString `db:"engine_type"`
    EngineVolume     sql.NullString `db:"engine_volume"`
    EnginePower      sql.NullInt64  `db:"engine_power"`
    Cylinders        sql.NullInt64  `db:"cylinders"`
    Transmission     sql.NullString `db:"transmission"`
    Drivetrain       sql.NullString `db:"drivetrain"`
    ModelID          sql.NullString `db:"model_id"`
    ModelName        sql.NullString `db:"model_name"`
    Make             sql.NullString `db:"make"`
    Year             sql.NullInt64  `db:"year"`
    BodyType         sql.NullString `db:"body_type"`
    Generation       sql.NullString `db:"generation"`
    WeightKg         sql.NullFloat64`db:"weight_kg"`
    SellerName       sql.NullString `db:"seller_name"`
    SellerRating     sql.NullFloat64`db:"seller_rating"`
    SellerSalesCount sql.NullInt64  `db:"seller_sales_count"`
    SellerIsBusiness sql.NullBool   `db:"seller_is_business"`
}

type ListingRepository interface {
    List(ctx context.Context, page, pageSize int, sortBy string) ([]Listing, error)
    Search(ctx context.Context, query string, page, pageSize int, sortBy string) ([]Listing, error)
    GetByID(ctx context.Context, listingID string) (*Listing, error)

    Create(ctx context.Context, l *Listing) error
    Update(ctx context.Context, l *Listing) error
    Delete(ctx context.Context, listingID string) error
}

type listingRepoDB struct {
    db *sql.DB
}

func NewListingRepository(db *sql.DB) ListingRepository {
    return &listingRepoDB{db: db}
}

func (r *listingRepoDB) List(ctx context.Context, page, pageSize int, sortBy string) ([]Listing, error) {
    offset := (page - 1) * pageSize
    orderClause := mapSortClause(sortBy)

    query := fmt.Sprintf(`
        SELECT
            listing_id,
            seller_id,
            description,
            posted_at,
            status,
            deal_type,
            price,
            car_id,
            mileage,
            owners_count,
            accidents_count,
            condition,
            color,
            config_id,
            engine_type,
            engine_volume,
            engine_power,
            cylinders,
            transmission,
            drivetrain,
            model_id,
            model_name,
            make,
            year,
            body_type,
            generation,
            weight_kg,
            seller_name,
            seller_rating,
            seller_sales_count,
            seller_is_business
        FROM listings
        %s
        LIMIT $1 OFFSET $2
    `, orderClause)

    rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
    if err != nil {
        return nil, fmt.Errorf("List query failed: %w", err)
    }
    defer rows.Close()

    var listings []Listing
    for rows.Next() {
        var l Listing
        if err := rows.Scan(
            &l.ListingID,
            &l.SellerID,
            &l.Description,
            &l.PostedAt,
            &l.Status,
            &l.DealType,
            &l.Price,
            &l.CarID,
            &l.Mileage,
            &l.OwnersCount,
            &l.AccidentsCount,
            &l.Condition,
            &l.Color,
            &l.ConfigID,
            &l.EngineType,
            &l.EngineVolume,
            &l.EnginePower,
            &l.Cylinders,
            &l.Transmission,
            &l.Drivetrain,
            &l.ModelID,
            &l.ModelName,
            &l.Make,
            &l.Year,
            &l.BodyType,
            &l.Generation,
            &l.WeightKg,
            &l.SellerName,
            &l.SellerRating,
            &l.SellerSalesCount,
            &l.SellerIsBusiness,
        ); err != nil {
            return nil, fmt.Errorf("List scan failed: %w", err)
        }
        listings = append(listings, l)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("List rows.Err: %w", err)
    }
    return listings, nil
}

func (r *listingRepoDB) Search(ctx context.Context, queryText string, page, pageSize int, sortBy string) ([]Listing, error) {
    offset := (page - 1) * pageSize
    orderClause := mapSortClause(sortBy)

    ilike := "%" + strings.ToLower(queryText) + "%"

    query := fmt.Sprintf(`
        SELECT
            listing_id,
            seller_id,
            description,
            posted_at,
            status,
            deal_type,
            price,
            car_id,
            mileage,
            owners_count,
            accidents_count,
            condition,
            color,
            config_id,
            engine_type,
            engine_volume,
            engine_power,
            cylinders,
            transmission,
            drivetrain,
            model_id,
            model_name,
            make,
            year,
            body_type,
            generation,
            weight_kg,
            seller_name,
            seller_rating,
            seller_sales_count,
            seller_is_business
        FROM listings
        WHERE
            LOWER(make) LIKE $1
            OR LOWER(model_name) LIKE $1
            OR LOWER(description) LIKE $1
        %s
        LIMIT $2 OFFSET $3
    `, orderClause)

    rows, err := r.db.QueryContext(ctx, query, ilike, pageSize, offset)
    if err != nil {
        return nil, fmt.Errorf("Search query failed: %w", err)
    }
    defer rows.Close()

    var listings []Listing
    for rows.Next() {
        var l Listing
        if err := rows.Scan(
            &l.ListingID,
            &l.SellerID,
            &l.Description,
            &l.PostedAt,
            &l.Status,
            &l.DealType,
            &l.Price,
            &l.CarID,
            &l.Mileage,
            &l.OwnersCount,
            &l.AccidentsCount,
            &l.Condition,
            &l.Color,
            &l.ConfigID,
            &l.EngineType,
            &l.EngineVolume,
            &l.EnginePower,
            &l.Cylinders,
            &l.Transmission,
            &l.Drivetrain,
            &l.ModelID,
            &l.ModelName,
            &l.Make,
            &l.Year,
            &l.BodyType,
            &l.Generation,
            &l.WeightKg,
            &l.SellerName,
            &l.SellerRating,
            &l.SellerSalesCount,
            &l.SellerIsBusiness,
        ); err != nil {
            return nil, fmt.Errorf("Search scan failed: %w", err)
        }
        listings = append(listings, l)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("Search rows.Err: %w", err)
    }
    return listings, nil
}

func (r *listingRepoDB) GetByID(ctx context.Context, listingID string) (*Listing, error) {
    query := `
        SELECT
            listing_id,
            seller_id,
            description,
            posted_at,
            status,
            deal_type,
            price,
            car_id,
            mileage,
            owners_count,
            accidents_count,
            condition,
            color,
            config_id,
            engine_type,
            engine_volume,
            engine_power,
            cylinders,
            transmission,
            drivetrain,
            model_id,
            model_name,
            make,
            year,
            body_type,
            generation,
            weight_kg,
            seller_name,
            seller_rating,
            seller_sales_count,
            seller_is_business
        FROM listings
        WHERE listing_id = $1
    `
    row := r.db.QueryRowContext(ctx, query, listingID)
    var l Listing
    if err := row.Scan(
        &l.ListingID,
        &l.SellerID,
        &l.Description,
        &l.PostedAt,
        &l.Status,
        &l.DealType,
        &l.Price,
        &l.CarID,
        &l.Mileage,
        &l.OwnersCount,
        &l.AccidentsCount,
        &l.Condition,
        &l.Color,
        &l.ConfigID,
        &l.EngineType,
        &l.EngineVolume,
        &l.EnginePower,
        &l.Cylinders,
        &l.Transmission,
        &l.Drivetrain,
        &l.ModelID,
        &l.ModelName,
        &l.Make,
        &l.Year,
        &l.BodyType,
        &l.Generation,
        &l.WeightKg,
        &l.SellerName,
        &l.SellerRating,
        &l.SellerSalesCount,
        &l.SellerIsBusiness,
    ); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("GetByID scan failed: %w", err)
    }
    return &l, nil
}

func (r *listingRepoDB) Create(ctx context.Context, l *Listing) error {
    _, err := r.db.ExecContext(ctx, `
        INSERT INTO listings (
          listing_id, seller_id, description, posted_at, status, deal_type,
          price, car_id, mileage, owners_count, accidents_count,
          condition, color, config_id, engine_type, engine_volume,
          engine_power, cylinders, transmission, drivetrain,
          model_id, model_name, make, year, body_type, generation,
          weight_kg, seller_name, seller_rating, seller_sales_count,
          seller_is_business
        ) VALUES (
          $1,$2,$3,$4,$5,$6,
          $7,$8,$9,$10,$11,
          $12,$13,$14,$15,$16,
          $17,$18,$19,$20,$21,$22,
          $23,$24,$25,$26,$27,$28,$29,$30,
          $31
        )
    `,
        l.ListingID,
        l.SellerID,
        l.Description, l.PostedAt, l.Status, l.DealType,
        l.Price, l.CarID, l.Mileage, l.OwnersCount, l.AccidentsCount,
        l.Condition, l.Color, l.ConfigID, l.EngineType, l.EngineVolume,
        l.EnginePower, l.Cylinders, l.Transmission, l.Drivetrain,
        l.ModelID, l.ModelName, l.Make, l.Year, l.BodyType, l.Generation,
        l.WeightKg, l.SellerName, l.SellerRating, l.SellerSalesCount,
        l.SellerIsBusiness,
    )
    if err != nil {
        return fmt.Errorf("Create listing failed: %w", err)
    }
    return nil
}

func (r *listingRepoDB) Update(ctx context.Context, l *Listing) error {
    _, err := r.db.ExecContext(ctx, `
        UPDATE listings SET
          seller_id          = $2,
          description        = $3,
          posted_at          = $4,
          status             = $5,
          deal_type          = $6,
          price              = $7,
          car_id             = $8,
          mileage            = $9,
          owners_count       = $10,
          accidents_count    = $11,
          condition          = $12,
          color              = $13,
          config_id          = $14,
          engine_type        = $15,
          engine_volume      = $16,
          engine_power       = $17,
          cylinders          = $18,
          transmission       = $19,
          drivetrain         = $20,
          model_id           = $21,
          model_name         = $22,
          make               = $23,
          year               = $24,
          body_type          = $25,
          generation         = $26,
          weight_kg          = $27,
          seller_name        = $28,
          seller_rating      = $29,
          seller_sales_count = $30,
          seller_is_business = $31
        WHERE listing_id = $1
    `,
        l.ListingID,
        l.SellerID,
        l.Description, l.PostedAt, l.Status, l.DealType,
        l.Price, l.CarID, l.Mileage, l.OwnersCount, l.AccidentsCount,
        l.Condition, l.Color, l.ConfigID, l.EngineType, l.EngineVolume,
        l.EnginePower, l.Cylinders, l.Transmission, l.Drivetrain,
        l.ModelID, l.ModelName, l.Make, l.Year, l.BodyType, l.Generation,
        l.WeightKg, l.SellerName, l.SellerRating, l.SellerSalesCount,
        l.SellerIsBusiness,
    )
    if err != nil {
        return fmt.Errorf("Update listing failed: %w", err)
    }
    return nil
}

func (r *listingRepoDB) Delete(ctx context.Context, listingID string) error {
    _, err := r.db.ExecContext(ctx,
        "DELETE FROM listings WHERE listing_id = $1",
        listingID,
    )
    if err != nil {
        return fmt.Errorf("Delete listing failed: %w", err)
    }
    return nil
}

func mapSortClause(sortBy string) string {
    switch sortBy {
    case "date_desc":
        return "ORDER BY posted_at DESC"
    case "date_asc":
        return "ORDER BY posted_at ASC"
    case "price_desc":
        return "ORDER BY price DESC"
    case "price_asc":
        return "ORDER BY price ASC"
    case "mileage_desc":
        return "ORDER BY mileage DESC"
    case "mileage_asc":
        return "ORDER BY mileage ASC"
    default:
        return "ORDER BY posted_at DESC"
    }
}