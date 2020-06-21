package server

import (
	"context"
	"net"
	"net/http"

	"github.com/xabi93/lana-test/internal/adding"
	"github.com/xabi93/lana-test/internal/calculating"
	"github.com/xabi93/lana-test/internal/creating"
	"github.com/xabi93/lana-test/internal/deleting"
	"github.com/xabi93/lana-test/pkg/log"

	"github.com/gorilla/mux"
)

// New initialize the server
func New(
	logger log.Logger,
	creating *creating.BasketCreator,
	adding *adding.BasketItemsAdder,
	deleting *deleting.BasketDeleter,
	calculating *calculating.TotalBasket,
) *Server {
	s := Server{
		logger:      logger,
		creating:    creating,
		adding:      adding,
		deleting:    deleting,
		calculating: calculating,
	}

	s.init()

	return &s
}

// Server all server necessary dependencies
type Server struct {
	router *mux.Router

	logger log.Logger

	creating    *creating.BasketCreator
	adding      *adding.BasketItemsAdder
	deleting    *deleting.BasketDeleter
	calculating *calculating.TotalBasket
}

func (s *Server) init() {
	s.initRoutes()
}

func (s *Server) initRoutes() {
	s.router = mux.NewRouter()

	s.router.HandleFunc("/baskets", s.handleBasketCreate()).Methods(http.MethodPost)
	s.router.HandleFunc("/baskets/{id}/items", s.handleAddItemBasket()).Methods(http.MethodPost)
	s.router.HandleFunc("/baskets/{id}/total", s.handleTotalBasket()).Methods(http.MethodGet)
	s.router.HandleFunc("/baskets/{id}", s.handleBasketDelete()).Methods(http.MethodDelete)

	s.router.Use(
		s.ContextMiddleware,
		s.LoggingMiddleware,
	)
}

func (s Server) Run(addr string) {
	if addr == "" {
		addr = "3000"
	}

	s.logger.Info(context.Background(), "Server running on: %s", addr)

	s.logger.Fatal(context.Background(), http.ListenAndServe(net.JoinHostPort("", addr), s.router))
}
