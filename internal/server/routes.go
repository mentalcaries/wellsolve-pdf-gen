package server

import (
	"fmt"
	"net/http"
)

var allowedOrigins = map[string]bool{
	"https://wellsolveable.com": true,
	"http:localhost:3000":       true,
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", s.handleReadiness)
	mux.HandleFunc("POST /pdf", s.generatePDF)

	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
		}

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleReadiness(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	resp := map[string]string{"message": "WellSolveAble PDF Service Online 🚀"}
	respondWithJSON(w, 200, resp)
}
