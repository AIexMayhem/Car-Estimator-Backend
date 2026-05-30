package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gorilla/mux"
	"github.com/nikita-itmo-gh-acc/car_estimator_api_gateway/config"
)

type IHandler interface {
	setupRoutes()
	setupgRPC(conn *grpc.ClientConn)
}

type Server struct {
	r *mux.Router
	port string
	logger *slog.Logger
	handlers map[string]IHandler
}

func NewServer(conf *config.Config, logger *slog.Logger) *Server {
	s := new(Server)
	s.r = mux.NewRouter()
	s.logger = logger
	s.port = conf.Port
	s.handlers = map[string]IHandler{}

	defer func(){
		if r := recover(); r != nil {
			s.logger.Error(r.(string))
		}
	}()

	s.logger.Info("Register services...")
	
	s.RegisterHandler("profile", &ProfileHandler{
		r: s.r.PathPrefix("/profile").Subrouter(),
		logger: s.logger, 
	}, MustConnect(conf.ProfileServiceAddr))
	
	s.RegisterHandler("feed", &FeedHandler{
		r: s.r.PathPrefix("/feed").Subrouter(),
		logger: s.logger, 
	}, MustConnect(conf.FeedServiceAddr))

	s.RegisterHandler("prediction", &PredictionHandler{
		r: s.r.PathPrefix("/prediction").Subrouter(),
		logger: s.logger, 
	}, MustConnect(conf.PredictionServiceAddr))

	s.logger.Info("Handlers registration completed!")

	return s
}

func MustConnect(addr string) *grpc.ClientConn {
	var (
		wait time.Duration = time.Second
		err error
		cc *grpc.ClientConn
	)

	for attempt := 0; attempt < 5; attempt++ {
		cc, err = grpc.NewClient(
			addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err == nil {
			break
		}

		time.Sleep(wait)
		wait = wait * 2
	}

	if err != nil {
		panic(fmt.Sprintf("can't connect to %s, error: %v", addr, err))
	}
	
	return cc
}

func (s *Server) RegisterHandler(name string, h IHandler, conn *grpc.ClientConn) {
	if _, ok := s.handlers[name]; ok {
		s.logger.Warn(fmt.Sprintf("attempt to recreate handler %s", name))
		return
	}

	h.setupgRPC(conn)
	s.handlers[name] = h
}

func (s *Server) Run() error {
	for _, handler := range s.handlers {
		handler.setupRoutes()
	}

	s.logger.Info("Running API Gateway server...", slog.String("port", s.port))
	return http.ListenAndServe(":" + s.port, s.r)
}
