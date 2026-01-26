package server

import (
	"encoding/json"
	"fmt"
	"io"
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
	mux.HandleFunc("POST /health", s.handleHealthCheck)

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

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(s.gotenbergEndpoint + "/health")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "PDF Service unreachable", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read health check response", err)
		return
	}

	var healthData map[string]interface{}
	if err := json.Unmarshal(body, &healthData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid health check response", err)
		return
	}

	respondWithJSON(w, http.StatusOK, healthData)
}
