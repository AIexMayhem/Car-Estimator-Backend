package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"log/slog"

	"github.com/gorilla/mux"
	feed "github.com/nikita-itmo-gh-acc/car_estimator_api_contracts/gen/feed_v1"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/domain"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/mappers"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/utils"
	"google.golang.org/grpc"
)

type FeedHandler struct {
	r      *mux.Router
	logger *slog.Logger
	client feed.FeedServiceClient
}

func (h *FeedHandler) setupgRPC(conn *grpc.ClientConn) {
	h.client = feed.NewFeedServiceClient(conn)
}

func (h *FeedHandler) setupRoutes() {
	h.r.HandleFunc("/listings", h.ListListings).Methods("GET")
	h.r.HandleFunc("/listings/search", h.SearchListings).Methods("GET")
	h.r.HandleFunc("/listings/{listingId}", h.GetListing).Methods("GET")
	h.r.HandleFunc("/listings", h.CreateListing).Methods("POST")
	h.r.HandleFunc("/listings/{listingId}", h.UpdateListing).Methods("PUT")
	h.r.HandleFunc("/listings/{listingId}", h.DeleteListing).Methods("DELETE")
	h.r.HandleFunc("/users/{userId}/favorites", h.AddToFavorites).Methods("POST")
}

func (h *FeedHandler) ListListings(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "ListListings"))

	q := r.URL.Query()
	pageNum, _ := strconv.Atoi(q.Get("page_number"))
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}
	sortByStr := strings.ToUpper(q.Get("sort_by"))
	sortBy := feed.SortBy_SORT_UNSPECIFIED
	if v, ok := feed.SortBy_value[sortByStr]; ok {
		sortBy = feed.SortBy(v)
	}

	grpcReq := &feed.ListListingsRequest{
		Page: &feed.PageRequest{
			PageNumber: int32(pageNum),
			PageSize:   int32(pageSize),
		},
		SortBy: sortBy,
	}

	log.Info("→ gRPC ListListings", slog.Any("req", grpcReq))
	grpcResp, err := h.client.ListListings(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "ListListings failed: ", err)
		return
	}

	out := domain.ListListingsResponse{
		PageMetadata: domain.PageResponseMetadata{
			TotalItems:  grpcResp.PageMetadata.TotalItems,
			TotalPages:  grpcResp.PageMetadata.TotalPages,
			CurrentPage: grpcResp.PageMetadata.CurrentPage,
		},
		Listings: make([]domain.CarListing, len(grpcResp.Listings)),
	}
	for i, c := range grpcResp.Listings {
		out.Listings[i] = *mappers.ToDomain(c)
	}

	utils.RenderJson(w, out)
}

func (h *FeedHandler) SearchListings(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "SearchListings"))

	q := r.URL.Query()
	query := q.Get("query")
	pageNum, _ := strconv.Atoi(q.Get("page_number"))
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}
	sortByStr := strings.ToUpper(q.Get("sort_by"))
	sortBy := feed.SortBy_SORT_UNSPECIFIED
	if v, ok := feed.SortBy_value[sortByStr]; ok {
		sortBy = feed.SortBy(v)
	}

	grpcReq := &feed.SearchListingsRequest{
		Query: query,
		Page: &feed.PageRequest{
			PageNumber: int32(pageNum),
			PageSize:   int32(pageSize),
		},
		SortBy: sortBy,
	}

	log.Info("→ gRPC SearchListings", slog.Any("req", grpcReq))
	grpcResp, err := h.client.SearchListings(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "SearchListings failed: ", err)
		return
	}

	out := domain.SearchListingsResponse{
		PageMetadata: domain.PageResponseMetadata{
			TotalItems:  grpcResp.PageMetadata.TotalItems,
			TotalPages:  grpcResp.PageMetadata.TotalPages,
			CurrentPage: grpcResp.PageMetadata.CurrentPage,
		},
		Listings: make([]domain.CarListing, len(grpcResp.Listings)),
	}
	for i, c := range grpcResp.Listings {
		out.Listings[i] = *mappers.ToDomain(c)
	}

	utils.RenderJson(w, out)
}

func (h *FeedHandler) GetListing(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "GetListing"))

	id := mux.Vars(r)["listingId"]
	grpcReq := &feed.GetListingRequest{ListingId: id}

	log.Info("→ gRPC GetListing", slog.String("id", id))
	grpcResp, err := h.client.GetListing(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "GetListing failed: ", err)
		return
	}

	out := domain.GetListingResponse{
		Listing: *mappers.ToDomain(grpcResp.Listing),
	}

	utils.RenderJson(w, out)
}

func (h *FeedHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "CreateListing"))

	var body domain.CreateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := &feed.CreateListingRequest{
		Listing: mappers.ToMessage(&body.Listing),
	}
	log.Info("→ gRPC CreateListing", slog.Any("req", grpcReq))

	grpcResp, err := h.client.CreateListing(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "CreateListing failed: ", err)
		return
	}

	out := domain.CreateListingResponse{
		Listing: *mappers.ToDomain(grpcResp.GetListing()),
	}
	utils.RenderJson(w, out)
}

func (h *FeedHandler) UpdateListing(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "UpdateListing"))

	listingId := mux.Vars(r)["listingId"]
	var body domain.UpdateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	body.Listing.ListingId = listingId

	grpcReq := &feed.UpdateListingRequest{
		Listing: mappers.ToMessage(&body.Listing),
	}

	log.Info("→ gRPC UpdateListing", slog.Any("req", grpcReq))
	grpcResp, err := h.client.UpdateListing(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "UpdateListing failed: ", err)
		return
	}

	out := domain.UpdateListingResponse{
		Listing: *mappers.ToDomain(grpcResp.GetListing()),
	}
	utils.RenderJson(w, out)
}

func (h *FeedHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "DeleteListing"))

	listingId := mux.Vars(r)["listingId"]
	log.Info("→ gRPC DeleteListing", slog.String("Id", listingId))

	grpcResp, err := h.client.DeleteListing(r.Context(), &feed.DeleteListingRequest{
		ListingId: listingId,
	})
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "DeleteListing failed: ", err)
		return
	}

	out := domain.DeleteListingResponse{Success: grpcResp.Success}
	utils.RenderJson(w, out)
}

func (h *FeedHandler) AddToFavorites(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(slog.String("op", "AddToFavorites"))

	userId := mux.Vars(r)["userId"]
	var body domain.AddToFavoritesRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	grpcReq := &feed.AddToFavoritesRequest{
		UserId:    userId,
		ListingId: body.ListingId,
	}

	log.Info("→ gRPC AddToFavorites", slog.Any("req", grpcReq))
	grpcResp, err := h.client.AddToFavorites(r.Context(), grpcReq)
	if err != nil {
		utils.HandleResponseErr(w, h.logger, "AddToFavorites failed: ", err)
		return
	}

	out := domain.AddToFavoritesResponse{Success: grpcResp.Success}
	utils.RenderJson(w, out)
}
