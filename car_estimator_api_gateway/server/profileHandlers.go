package server

import (
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	profile "github.com/nikita-itmo-gh-acc/car_estimator_api_contracts/gen/profile_v1"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/domain"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/utils"
)

type ProfileHandler struct {
	r *mux.Router
	logger *slog.Logger
	client profile.ProfileServiceClient
}

func (h *ProfileHandler) setupgRPC(conn *grpc.ClientConn) {
	h.client = profile.NewProfileServiceClient(conn)
}

func (h *ProfileHandler) setupRoutes() {
	h.r.HandleFunc("/login", h.LoginHandler).Methods("POST")
	h.r.HandleFunc("/logout", h.LogoutHandler).Methods("DELETE")
	h.r.HandleFunc("/register", h.RegisterHandler).Methods("POST")
	h.r.HandleFunc("/unregister", h.UnregisterHandler).Methods("DELETE")
	h.r.HandleFunc("/users/{userId}", h.GetUserHandler).Methods("GET")
	h.r.HandleFunc("/refresh", h.RefreshHandler).Methods("POST")
}

func (h *ProfileHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	body := &domain.LoginRequest{}
	if err := utils.ParseJson(r.Body, body); err != nil {
		http.Error(w, "login request body decode failed", http.StatusBadRequest)
		return
	}

	ipAddr := utils.GetClientIp(r) 
	userAgent := r.UserAgent()

	log := h.logger.With(
		slog.String("operation", "login"),
		slog.String("email", body.Email),
		slog.String("source", ipAddr + " " + userAgent),
	)

	log.Info("Attempt to login. Calling profile gRPC service...")

	response, err := h.client.Login(r.Context(), &profile.LoginRequest{
		Email: body.Email,
		Password: body.Password,
		Source: &profile.SourceData{
			Ip: ipAddr,
			UserAgent: userAgent,
		},
	})

	if err != nil {
		utils.HandleResponseErr(w, h.logger, "login failed - ", err)
		return
	}

	log.Info("Successfully logged in!")
	id, _ := uuid.Parse(response.GetUserId().Value)

	utils.RenderJson(w, domain.LoginResponse{
		UserId: id,
		TokenResponse: domain.TokenResponse{
			Access: response.GetTokens().AccessToken,
			Refresh: response.GetTokens().RefreshToken,
		},
	})
}

func (h *ProfileHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	source := utils.GetClientIp(r) + " " + r.UserAgent()

	log := h.logger.With(
		slog.String("opetation", "logout"),
		slog.String("source", source),
	)

	log.Info("Attempt to logout. Calling profile gRPC service...")

	md := metadata.Pairs(
		"refreshToken", r.Header.Get("refreshToken"),
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)
	if _, err := h.client.Logout(ctx, &emptypb.Empty{}); err != nil {
		utils.HandleResponseErr(w, h.logger, "logout failed - ", err)
		return
	}

	log.Info("Successfully logged out!")
}

func (h *ProfileHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(
		slog.String("operation", "get user"),
	)

	log.Info("Start user search process. Calling profile gRPC service...")

	userId, ok := mux.Vars(r)["userId"]
	if !ok {
		http.Error(w, "param userId is missing", http.StatusBadRequest)
		return
	}

	_, err := uuid.Parse(userId)
	if err != nil {
		http.Error(w, "param userId has wrong format", http.StatusBadRequest)
		return
	}

	response, err := h.client.GetUser(r.Context(), &profile.UserRequest{
		UserId: &profile.UUID{
			Value: userId,
		},
	})

	if err != nil {
		utils.HandleResponseErr(w, h.logger, "get user op failed - ", err)
		return
	}

	log.Info("Successfully retreived user!",
		slog.String("user ID", userId),
	)

	utils.RenderJson(w, domain.UserResponse{
		Id: uuid.MustParse(response.GetUserId().GetValue()),
		FullName: response.GetFullname(),
		Email: response.GetEmail(),
		Phone: response.GetPhone(),
		BirthDate: time.Unix(response.GetBirthdate(), 0),
		RegisterDate: time.Unix(response.GetRegisterdate(), 0),
	})
}

func (h *ProfileHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(
		slog.String("operation", "refresh tokens"),
	)

	log.Info("Start refreshing process. Calling profile gRPC service...")

	refreshToken := r.Header.Get("RefreshToken")
	ipAddr := utils.GetClientIp(r) 
	userAgent := r.UserAgent()

	md := metadata.Pairs(
		"refreshToken", refreshToken,
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)
	response, err := h.client.Refresh(ctx, &profile.SourceData{
		Ip: ipAddr,
		UserAgent: userAgent,
	})

	if err != nil {
		utils.HandleResponseErr(w, h.logger, "refresh op failed - ", err)
		return
	}

	utils.RenderJson(w, domain.TokenResponse{
		Access: response.GetAccessToken(),
		Refresh: response.GetRefreshToken(),
	})
}

func (h *ProfileHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	body := &domain.RegisterRequest{}
	if err := utils.ParseJson(r.Body, body); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "register request body decode failed", http.StatusBadRequest)
		return
	}

	log := h.logger.With(
		slog.String("operation", "register"),
		slog.String("username", body.FullName),
		slog.String("email", body.Email),
	)

	bd, err := time.Parse("2006-01-02", body.BirthDate)
	if err != nil {
		http.Error(w, "birthDate has wrong format", http.StatusBadRequest)
		return
	}

	log.Info("Start register process. Calling profile gRPC service...")

	response, err := h.client.Register(r.Context(), &profile.RegisterRequest{
		Fullname: body.FullName,
		Email: body.Email,
		Phone: body.Phone,
		Password: body.Password,
		Birthdate: bd.Unix(),
	})

	if err != nil {
		utils.HandleResponseErr(w, h.logger, "register failed - ", err)
		return
	}

	id, _ := uuid.Parse(response.GetUserId().Value)
	log.Info("Successfully registered!")

	utils.RenderJson(w, domain.RegisterResponse{
		UserId: id,
	})
}

func (h *ProfileHandler) UnregisterHandler(w http.ResponseWriter, r *http.Request) {
	source := utils.GetClientIp(r) + " " + r.UserAgent()

	log := h.logger.With(
		slog.String("opetation", "logout"),
		slog.String("source", source),
	)

	log.Info("Start	unregister process. Calling profile gRPC service...")

	md := metadata.Pairs(
		"refreshToken", r.Header.Get("refreshToken"),
	)

	ctx := metadata.NewOutgoingContext(r.Context(), md)
	if _, err := h.client.Unregister(ctx, &emptypb.Empty{}); err != nil {
		utils.HandleResponseErr(w, h.logger, "unregister failed - ", err)
		return
	}

	log.Info("Successfully unregistered!")
}
