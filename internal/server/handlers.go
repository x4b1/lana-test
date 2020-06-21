package server

import (
	"context"
	"encoding/json"
	"net/http"

	checkout "github.com/xabi93/lana-test/internal"
	"github.com/xabi93/lana-test/pkg/errors"

	"github.com/gorilla/mux"
)

func (s Server) handleBasketCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := s.creating.Create(r.Context())
		if err != nil {
			s.responseError(r.Context(), w, err)
			return
		}

		s.respond(w, http.StatusCreated, basketCreated{ID: string(b.ID)})
	}
}

func (s Server) handleAddItemBasket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil || r.Body == http.NoBody {
			s.responseError(r.Context(), w, errors.NewWrongInput("Empty body"))
			return
		}
		var request struct {
			Product string `json:"product"`
		}

		if err := decodeRequest(r, &request); err != nil {
			s.responseError(r.Context(), w, err)
			return
		}

		id := checkout.BasketID(mux.Vars(r)["id"])

		if request.Product == "" {
			s.responseError(r.Context(), w, errors.NewWrongInput("Empty product"))
			return
		}

		if err := s.adding.Add(r.Context(), id, checkout.ProductCode(request.Product)); err != nil {
			s.responseError(r.Context(), w, err)
			return
		}

		s.respond(w, http.StatusCreated, nil)
	}
}

func (s Server) handleBasketDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := checkout.BasketID(mux.Vars(r)["id"])

		if err := s.deleting.Delete(r.Context(), id); err != nil {
			s.responseError(r.Context(), w, err)
			return
		}

		s.respond(w, http.StatusNoContent, nil)
	}
}

func (s Server) handleTotalBasket() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := checkout.BasketID(mux.Vars(r)["id"])

		total, err := s.calculating.Total(r.Context(), id)
		if err != nil {
			s.responseError(r.Context(), w, err)
			return
		}

		s.respond(w, http.StatusOK, basketTotal{Total: total.String()})
	}
}

func (s Server) responseError(ctx context.Context, w http.ResponseWriter, err error) {
	var code int
	switch {
	case errors.IsWrongInput(err):
		code = http.StatusBadRequest
	case errors.IsNotFound(err):
		code = http.StatusNotFound
	default:
		code = http.StatusInternalServerError
	}

	s.logger.Error(ctx, err)

	s.respond(w, code, errorResponse{Message: err.Error()})
}

func (s Server) respond(w http.ResponseWriter, code int, msg interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(msg)
}

func decodeRequest(r *http.Request, payload interface{}) error {
	return json.NewDecoder(r.Body).Decode(payload)
}
