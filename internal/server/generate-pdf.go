package server

import (
	"fmt"
	"net/http"
)

func(s *Server) generatePDF(w http.ResponseWriter, r *http.Request){
  fmt.Println("Generating PDF")
}