package grpc_server

import (
    "context"
    "fmt"
    "net"

    "log/slog"

    feedcontract "github.com/nikita-itmo-gh-acc/car_estimator_api_contracts/gen/feed_v1"
    "github.com/icecarti/feed_service/services"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

type GRPCServer struct {
    grpcServer *grpc.Server
    port       string
}

func NewGRPCServer(port string, feedSvc services.FeedService, logger *slog.Logger) *GRPCServer {
    unaryLog := grpc.UnaryInterceptor(func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        logger.Info("gRPC request", "method", info.FullMethod, "request", req)
        resp, err := handler(ctx, req)
        if err != nil {
            logger.Error("gRPC error", "method", info.FullMethod, "error", err)
        }
        return resp, err
    })

    s := grpc.NewServer(unaryLog)
    feedcontract.RegisterFeedServiceServer(s, &feedGRPCWrapper{
        UnimplementedFeedServiceServer: feedcontract.UnimplementedFeedServiceServer{},
        svc:                            feedSvc,
    })
    reflection.Register(s)
    return &GRPCServer{grpcServer: s, port: port}
}

func (s *GRPCServer) Start() error {
    lis, err := net.Listen("tcp", s.port)
    if err != nil {
        return fmt.Errorf("failed to listen on %s: %w", s.port, err)
    }
    fmt.Printf("gRPC server listening on %s\n", s.port)
    return s.grpcServer.Serve(lis)
}

func (s *GRPCServer) Stop() {
    s.grpcServer.GracefulStop()
}

type feedGRPCWrapper struct {
    feedcontract.UnimplementedFeedServiceServer
    svc services.FeedService
}

func (w *feedGRPCWrapper) ListListings(
    ctx context.Context,
    req *feedcontract.ListListingsRequest,
) (*feedcontract.ListListingsResponse, error) {
    return w.svc.ListListings(ctx, req)
}

func (w *feedGRPCWrapper) GetListing(
    ctx context.Context,
    req *feedcontract.GetListingRequest,
) (*feedcontract.GetListingResponse, error) {
    return w.svc.GetListing(ctx, req)
}

func (w *feedGRPCWrapper) SearchListings(
    ctx context.Context,
    req *feedcontract.SearchListingsRequest,
) (*feedcontract.SearchListingsResponse, error) {
    return w.svc.SearchListings(ctx, req)
}

func (w *feedGRPCWrapper) CreateListing(
    ctx context.Context,
    req *feedcontract.CreateListingRequest,
) (*feedcontract.CreateListingResponse, error) {
    return w.svc.CreateListing(ctx, req)
}

func (w *feedGRPCWrapper) UpdateListing(
    ctx context.Context,
    req *feedcontract.UpdateListingRequest,
) (*feedcontract.UpdateListingResponse, error) {
    return w.svc.UpdateListing(ctx, req)
}

func (w *feedGRPCWrapper) DeleteListing(
    ctx context.Context,
    req *feedcontract.DeleteListingRequest,
) (*feedcontract.DeleteListingResponse, error) {
    return w.svc.DeleteListing(ctx, req)
}

func (w *feedGRPCWrapper) AddToFavorites(
    ctx context.Context,
    req *feedcontract.AddToFavoritesRequest,
) (*feedcontract.AddToFavoritesResponse, error) {
    return w.svc.AddToFavorites(ctx, req)
}