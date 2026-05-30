package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/icecarti/feed_service/database"
	feedcontract "github.com/nikita-itmo-gh-acc/car_estimator_api_contracts/gen/feed_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FeedService interface {
    ListListings(ctx context.Context, req *feedcontract.ListListingsRequest) (*feedcontract.ListListingsResponse, error)
    SearchListings(ctx context.Context, req *feedcontract.SearchListingsRequest) (*feedcontract.SearchListingsResponse, error)
    GetListing(ctx context.Context, req *feedcontract.GetListingRequest) (*feedcontract.GetListingResponse, error)
    CreateListing(ctx context.Context, req *feedcontract.CreateListingRequest) (*feedcontract.CreateListingResponse, error)
    UpdateListing(ctx context.Context, req *feedcontract.UpdateListingRequest) (*feedcontract.UpdateListingResponse, error)
    DeleteListing(ctx context.Context, req *feedcontract.DeleteListingRequest) (*feedcontract.DeleteListingResponse, error)
    AddToFavorites(ctx context.Context, req *feedcontract.AddToFavoritesRequest) (*feedcontract.AddToFavoritesResponse, error)
}

type feedServiceImpl struct {
    listingRepo  database.ListingRepository
    favoriteRepo database.FavoriteRepository
}

func NewFeedService(lr database.ListingRepository, fr database.FavoriteRepository) FeedService {
    return &feedServiceImpl{listingRepo: lr, favoriteRepo: fr}
}

func (s *feedServiceImpl) ListListings(ctx context.Context, req *feedcontract.ListListingsRequest) (*feedcontract.ListListingsResponse, error) {
    page := int(req.GetPage().GetPageNumber())
    pageSize := int(req.GetPage().GetPageSize())
    sortBy := mapSortBy(req.GetSortBy())

    listingsDB, err := s.listingRepo.List(ctx, page, pageSize, sortBy)
    if err != nil {
        return nil, fmt.Errorf("ListingRepository.List: %w", err)
    }

    protoListings := make([]*feedcontract.CarListing, 0, len(listingsDB))
    for _, dbL := range listingsDB {
        protoListings = append(protoListings, convertDBToProto(&dbL))
    }

    pageMetadata := &feedcontract.PageResponseMetadata{
        TotalItems:  int32(len(protoListings)),
        TotalPages:  1,
        CurrentPage: req.GetPage().GetPageNumber(),
    }

    return &feedcontract.ListListingsResponse{
        Listings:     protoListings,
        PageMetadata: pageMetadata,
    }, nil
}

func (s *feedServiceImpl) SearchListings(ctx context.Context, req *feedcontract.SearchListingsRequest) (*feedcontract.SearchListingsResponse, error) {
    page := int(req.GetPage().GetPageNumber())
    pageSize := int(req.GetPage().GetPageSize())
    sortBy := mapSortBy(req.GetSortBy())

    listingsDB, err := s.listingRepo.Search(ctx, req.GetQuery(), page, pageSize, sortBy)
    if err != nil {
        return nil, fmt.Errorf("ListingRepository.Search: %w", err)
    }

    protoListings := make([]*feedcontract.CarListing, 0, len(listingsDB))
    for _, dbL := range listingsDB {
        protoListings = append(protoListings, convertDBToProto(&dbL))
    }

    pageMetadata := &feedcontract.PageResponseMetadata{
        TotalItems:  int32(len(protoListings)),
        TotalPages:  1,
        CurrentPage: req.GetPage().GetPageNumber(),
    }

    return &feedcontract.SearchListingsResponse{
        Listings:     protoListings,
        PageMetadata: pageMetadata,
    }, nil
}

func (s *feedServiceImpl) GetListing(ctx context.Context, req *feedcontract.GetListingRequest) (*feedcontract.GetListingResponse, error) {
    dbL, err := s.listingRepo.GetByID(ctx, req.GetListingId())
    if err != nil {
        return nil, fmt.Errorf("ListingRepository.GetByID: %w", err)
    }
    if dbL == nil {
        return nil, fmt.Errorf("listing not found: %s", req.GetListingId())
    }
    return &feedcontract.GetListingResponse{
        Listing: convertDBToProto(dbL),
    }, nil
}

func (s *feedServiceImpl) CreateListing(ctx context.Context, req *feedcontract.CreateListingRequest) (*feedcontract.CreateListingResponse, error) {
    dbL := convertProtoToDB(req.GetListing())
    dbL.ListingID = uuid.NewString()
    dbL.PostedAt = time.Now()

    if err := s.listingRepo.Create(ctx, dbL); err != nil {
        return nil, fmt.Errorf("ListingRepository.Create: %w", err)
    }
    return &feedcontract.CreateListingResponse{
        Listing: convertDBToProto(dbL),
    }, nil
}

func (s *feedServiceImpl) UpdateListing(ctx context.Context, req *feedcontract.UpdateListingRequest) (*feedcontract.UpdateListingResponse, error) {
    dbL := convertProtoToDB(req.GetListing())
    if err := s.listingRepo.Update(ctx, dbL); err != nil {
        return nil, fmt.Errorf("ListingRepository.Update: %w", err)
    }
    return &feedcontract.UpdateListingResponse{
        Listing: convertDBToProto(dbL),
    }, nil
}

func (s *feedServiceImpl) DeleteListing(ctx context.Context, req *feedcontract.DeleteListingRequest) (*feedcontract.DeleteListingResponse, error) {
    if err := s.listingRepo.Delete(ctx, req.GetListingId()); err != nil {
        return &feedcontract.DeleteListingResponse{Success: false}, fmt.Errorf("ListingRepository.Delete: %w", err)
    }
    return &feedcontract.DeleteListingResponse{Success: true}, nil
}

func (s *feedServiceImpl) AddToFavorites(ctx context.Context, req *feedcontract.AddToFavoritesRequest) (*feedcontract.AddToFavoritesResponse, error) {
    if err := s.favoriteRepo.AddToFavorites(ctx, req.GetUserId(), req.GetListingId()); err != nil {
        return &feedcontract.AddToFavoritesResponse{Success: false}, fmt.Errorf("FavoriteRepository.AddToFavorites: %w", err)
    }
    return &feedcontract.AddToFavoritesResponse{Success: true}, nil
}

func mapSortBy(s feedcontract.SortBy) string {
    switch s {
    case feedcontract.SortBy_SORT_DATE_DESC:
        return "date_desc"
    case feedcontract.SortBy_SORT_DATE_ASC:
        return "date_asc"
    case feedcontract.SortBy_SORT_PRICE_DESC:
        return "price_desc"
    case feedcontract.SortBy_SORT_PRICE_ASC:
        return "price_asc"
    case feedcontract.SortBy_SORT_MILEAGE_DESC:
        return "mileage_desc"
    case feedcontract.SortBy_SORT_MILEAGE_ASC:
        return "mileage_asc"
    default:
        return ""
    }
}

func convertDBToProto(d *database.Listing) *feedcontract.CarListing {
    return &feedcontract.CarListing{
        ListingId:        d.ListingID,
        SellerId:         d.SellerID,
        Description:      d.Description.String,
        PostedAt:         timestamppb.New(d.PostedAt),
        Status:           d.Status.String,
        DealType:         d.DealType.String,
        Price:            d.Price.Float64,
        CarId:            d.CarID.String,
        Mileage:          int32(d.Mileage.Int64),
        OwnersCount:      int32(d.OwnersCount.Int64),
        AccidentsCount:   int32(d.AccidentsCount.Int64),
        Condition:        d.Condition.String,
        Color:            d.Color.String,
        ConfigId:         d.ConfigID.String,
        EngineType:       d.EngineType.String,
        EngineVolume:     d.EngineVolume.String,
        EnginePower:      int32(d.EnginePower.Int64),
        Cylinders:        int32(d.Cylinders.Int64),
        Transmission:     d.Transmission.String,
        Drivetrain:       d.Drivetrain.String,
        ModelId:          d.ModelID.String,
        ModelName:        d.ModelName.String,
        Make:             d.Make.String,
        Year:             int32(d.Year.Int64),
        BodyType:         d.BodyType.String,
        Generation:       d.Generation.String,
        WeightKg:         d.WeightKg.Float64,
        SellerName:       d.SellerName.String,
        SellerRating:     d.SellerRating.Float64,
        SellerSalesCount: int32(d.SellerSalesCount.Int64),
        SellerIsBusiness: d.SellerIsBusiness.Bool,
    }
}

func convertProtoToDB(p *feedcontract.CarListing) *database.Listing {
    return &database.Listing{
        ListingID:        p.GetListingId(),
        SellerID:         p.GetSellerId(),
        Description:      sql.NullString{String: p.GetDescription(), Valid: p.GetDescription() != ""},
        PostedAt:         p.GetPostedAt().AsTime(),
        Status:           sql.NullString{String: p.GetStatus(), Valid: p.GetStatus() != ""},
        DealType:         sql.NullString{String: p.GetDealType(), Valid: p.GetDealType() != ""},
        Price:            sql.NullFloat64{Float64: p.GetPrice(), Valid: p.GetPrice() != 0},
        CarID:            sql.NullString{String: p.GetCarId(), Valid: p.GetCarId() != ""},
        Mileage:          sql.NullInt64{Int64: int64(p.GetMileage()), Valid: p.GetMileage() != 0},
        OwnersCount:      sql.NullInt64{Int64: int64(p.GetOwnersCount()), Valid: p.GetOwnersCount() != 0},
        AccidentsCount:   sql.NullInt64{Int64: int64(p.GetAccidentsCount()), Valid: p.GetAccidentsCount() != 0},
        Condition:        sql.NullString{String: p.GetCondition(), Valid: p.GetCondition() != ""},
        Color:            sql.NullString{String: p.GetColor(), Valid: p.GetColor() != ""},
        ConfigID:         sql.NullString{String: p.GetConfigId(), Valid: p.GetConfigId() != ""},
        EngineType:       sql.NullString{String: p.GetEngineType(), Valid: p.GetEngineType() != ""},
        EngineVolume:     sql.NullString{String: p.GetEngineVolume(), Valid: p.GetEngineVolume() != ""},
        EnginePower:      sql.NullInt64{Int64: int64(p.GetEnginePower()), Valid: p.GetEnginePower() != 0},
        Cylinders:        sql.NullInt64{Int64: int64(p.GetCylinders()), Valid: p.GetCylinders() != 0},
        Transmission:     sql.NullString{String: p.GetTransmission(), Valid: p.GetTransmission() != ""},
        Drivetrain:       sql.NullString{String: p.GetDrivetrain(), Valid: p.GetDrivetrain() != ""},
        ModelID:          sql.NullString{String: p.GetModelId(), Valid: p.GetModelId() != ""},
        ModelName:        sql.NullString{String: p.GetModelName(), Valid: p.GetModelName() != ""},
        Make:             sql.NullString{String: p.GetMake(), Valid: p.GetMake() != ""},
        Year:             sql.NullInt64{Int64: int64(p.GetYear()), Valid: p.GetYear() != 0},
        BodyType:         sql.NullString{String: p.GetBodyType(), Valid: p.GetBodyType() != ""},
        Generation:       sql.NullString{String: p.GetGeneration(), Valid: p.GetGeneration() != ""},
        WeightKg:         sql.NullFloat64{Float64: p.GetWeightKg(), Valid: p.GetWeightKg() != 0},
        SellerName:       sql.NullString{String: p.GetSellerName(), Valid: p.GetSellerName() != ""},
        SellerRating:     sql.NullFloat64{Float64: p.GetSellerRating(), Valid: p.GetSellerRating() != 0},
        SellerSalesCount: sql.NullInt64{Int64: int64(p.GetSellerSalesCount()), Valid: p.GetSellerSalesCount() != 0},
        SellerIsBusiness: sql.NullBool{Bool: p.GetSellerIsBusiness(), Valid: true},
    }
}
