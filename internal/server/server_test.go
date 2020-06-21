package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xabi93/lana-test/internal/adding"
	"github.com/xabi93/lana-test/internal/calculating"
	"github.com/xabi93/lana-test/internal/creating"
	"github.com/xabi93/lana-test/internal/deleting"
	"github.com/xabi93/lana-test/internal/storage/memory"

	"github.com/stretchr/testify/suite"
)

type logger struct{}

func (l *logger) Error(ctx context.Context, err error) {}

func (m *logger) Fatal(ctx context.Context, err error) {}

func (m *logger) Info(ctx context.Context, msg string, args ...interface{}) {}

type ServerSuite struct {
	suite.Suite

	server *Server
}

func (s *ServerSuite) SetupTest() {
	db := memory.NewDB()
	baskets := memory.NewBasketRepository(db)
	products := memory.NewProductRepository(db)
	discounts := memory.NewDiscountRepository(db)

	s.server = New(
		new(logger),
		creating.NewBasketCreator(baskets),
		adding.NewBasketItemsAdder(baskets, products),
		deleting.NewBasketDeleter(baskets),
		calculating.NewTotalBasket(baskets, products, discounts),
	)
}

func (s *ServerSuite) TestCreateBasketReturnsCreated() {
	response := requestCreateBasket(s.server)

	s.Equal(http.StatusCreated, response.Code)
}
func (s *ServerSuite) TestAddItemWithNoProductReturnsBadRequest() {
	response := requestAddItem(s.server, "f1535292-a229-438a-b1a8-fcb05ed51435", "")

	s.Equal(http.StatusBadRequest, response.Code)
}

func (s *ServerSuite) TestAddItemToNotExistingBasketReturnsNotFound() {
	response := requestAddItem(s.server, "f1535292-a229-438a-b1a8-fcb05ed51435", "PEN")

	s.Equal(http.StatusNotFound, response.Code)
}

func (s *ServerSuite) TestAddItemProductNotExistsReturnsNotFound() {
	var basketResp basketCreated
	decodeResponse(requestCreateBasket(s.server), &basketResp)

	response := requestAddItem(s.server, basketResp.ID, "DRESS")

	s.Equal(http.StatusNotFound, response.Code)
}

func (s *ServerSuite) TestAddItemToExistingBasketReturnsCreated() {
	var basketResp basketCreated
	decodeResponse(requestCreateBasket(s.server), &basketResp)

	response := requestAddItem(s.server, basketResp.ID, "PEN")
	s.Equal(http.StatusCreated, response.Code)
}

func (s *ServerSuite) TestGetTotalNotExistingBasket() {
	response := requestGetTotal(s.server, "82d7c7f7-af4c-416a-8c16-879442f500c3")

	s.Equal(http.StatusNotFound, response.Code)
}

func (s *ServerSuite) TestGetTotalSucess() {
	type testCase struct {
		items    []string
		expected string
	}

	cases := map[string]testCase{
		"no discounts": {
			items:    []string{"PEN", "TSHIRT", "MUG"},
			expected: "32.50\u20ac",
		},
		"buy 2 pay 1": {
			items:    []string{"PEN", "TSHIRT", "PEN"},
			expected: "25.00\u20ac",
		},
		"t-shirt bulk purchase": {
			items:    []string{"TSHIRT", "TSHIRT", "TSHIRT", "PEN", "TSHIRT"},
			expected: "65.00\u20ac",
		},
		"both discounts": {
			items:    []string{"PEN", "TSHIRT", "PEN", "PEN", "MUG", "TSHIRT", "TSHIRT"},
			expected: "62.50\u20ac",
		},
	}

	for t, tc := range cases {
		s.Run(t, func() {
			var basketResp basketCreated
			decodeResponse(requestCreateBasket(s.server), &basketResp)

			for _, product := range tc.items {
				requestAddItem(s.server, basketResp.ID, product)
			}

			response := requestGetTotal(s.server, basketResp.ID)

			var totalResp basketTotal
			decodeResponse(response, &totalResp)

			s.Equal(http.StatusOK, response.Code)
			s.Equal(tc.expected, totalResp.Total)
		})
	}
}

func (s *ServerSuite) TestDeleteNotExistingBasketReturnsNotFound() {
	response := requestDeleteBasket(s.server, "541b22c9-b17c-4c34-a93d-92cc7118c97d")

	s.Equal(http.StatusNotFound, response.Code)
}

func (s *ServerSuite) TestDeleteExistingBasketDeletes() {
	var basketResp basketCreated
	decodeResponse(requestCreateBasket(s.server), &basketResp)

	response := requestDeleteBasket(s.server, basketResp.ID)
	s.Equal(http.StatusNoContent, response.Code)

	response = requestGetTotal(s.server, basketResp.ID)
	s.Equal(http.StatusNotFound, response.Code)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

func executeRequest(req *http.Request, server *Server) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)

	return recorder
}

func requestCreateBasket(s *Server) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, "/baskets", nil)

	return executeRequest(req, s)
}

func requestAddItem(s *Server, basket string, product string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/baskets/%s/items", basket),
		strings.NewReader(fmt.Sprintf(`{"product":"%s"}`, product)),
	)
	return executeRequest(req, s)
}

func requestGetTotal(s *Server, basket string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/baskets/%s/total", basket),
		nil,
	)
	return executeRequest(req, s)
}

func requestDeleteBasket(s *Server, basket string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/baskets/%s", basket), nil)

	return executeRequest(req, s)
}

func decodeResponse(r *httptest.ResponseRecorder, payload interface{}) error {
	return json.NewDecoder(r.Body).Decode(payload)
}
