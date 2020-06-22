package server

import (
	"context"
	"net/http"
)

//LoggingMiddleware logs the url path of the request
func (s *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(r.Context(), r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

//ContextMiddleware adds info of the request to the context
func (s *Server) ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ctx = context.WithValue(ctx, contextKeyEndpoint, r.URL.RequestURI())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
