package server

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"google.golang.org/grpc"

	"github.com/gorilla/mux"
	model "github.com/nikita-itmo-gh-acc/car_estimator_api_contracts/gen/prediction_v1/go"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/domain"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/utils"
)

type PredictionHandler struct {
	r *mux.Router
	logger *slog.Logger
	client model.PredictionServiceClient
}

func (h *PredictionHandler) setupgRPC(conn *grpc.ClientConn) {
	h.client = model.NewPredictionServiceClient(conn)
}

func (h *PredictionHandler) setupRoutes() {
	h.r.HandleFunc("", h.PredictionHandler).Methods("POST")
	h.r.HandleFunc("/images/{make}/{model}/{year}", h.GetImagesHandler).Methods("GET")
}

func (h *PredictionHandler) GetImagesHandler(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(
		slog.String("operation", "get car images"),
	)

	carInfo := map[string]string{
		"make": "",
		"model": "",
		"year": "",
	}

	for k := range carInfo {
		param, ok := mux.Vars(r)[k]
		if !ok {
			http.Error(w, fmt.Sprintf("param %s is missing", k), http.StatusBadRequest)
			return
		}
		carInfo[k] = param
	}

	log.Info(
		"Try to retreive car images by given params...",
		slog.Any("params", carInfo),
	)

	year, err := strconv.Atoi(carInfo["year"])
	if err != nil {
		http.Error(w, "year param has wrong format", http.StatusBadRequest)
		return
	}

	response, err := h.client.GetImages(r.Context(), &model.ImagesRequest{
		Make: carInfo["make"],
		Model: carInfo["model"],
		Year: int32(year),
	})

	if err != nil {
		utils.HandleResponseErr(w, h.logger, "image retreive operation failed - ", err)
		return
	}

	log.Info(
		"Successfully received car images!",
	)

	utils.RenderJson(w, domain.ImageResponse{
		Urls: response.PhotoUrls,
	})
}

func (h *PredictionHandler) PredictionHandler(w http.ResponseWriter, r *http.Request) {
	log := h.logger.With(
		slog.String("operation", "get car price prediction"),
	)

	params := new(domain.PredictionRequest)
	if err := utils.ParseJson(r.Body, params); err != nil {
		http.Error(w, "failed to decode prediction request body: " + err.Error(), http.StatusBadRequest)
		return
	}

	log.Info(
		"Send given params to the prediction service...",
		slog.Any("params", params),
	)

	response, err := h.client.Predict(r.Context(), &model.PredictRequest{
		Make: params.Make,
		Model: params.Model,
		Year: int32(params.Year),
		Hp: int32(params.Hp),
		Body: params.Body,
		Yearsell: int32(params.YearSell),
		Odometer: int32(params.Odometer),
		Color: params.Color,
	})

	if err != nil {
		utils.HandleResponseErr(w, h.logger, "prediction operation failed - ", err)
		return
	}

	imageB64 := base64.StdEncoding.EncodeToString(response.GraphPng)
	
	log.Info(
		"Successfully received prediction from the model!",
	)

	utils.RenderJson(w, &domain.PredictionResponse{
		Price: int(response.GetPrice()),
		SellCount: int(response.GetSellCount()),
		Urls: response.GetPhotoUrls(),
		GraphImg: imageB64,
	})
}
