package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type PDFRequest struct {
	HTML     string `json:"html"`
	Filename string `json:"filename,omitempty"`
}

func (s *Server) generatePDF(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	req := PDFRequest{}

	err := decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not decode request", err)
		return
	}

	defer r.Body.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	htmlPart, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create form data", err)
		return
	}

	_, err = htmlPart.Write([]byte(req.HTML))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to write HTML", err)
		return
	}

	err = writer.Close()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to finalize form", err)
		return
	}

	resp, err := http.Post(
		s.gotenbergEndpoint+"/forms/chromium/convert/html",
		writer.FormDataContentType(),
		&body,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate PDF", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respondWithError(w, resp.StatusCode, "PDF generation failed", nil)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=document.pdf")

	// Stream PDF back to client
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error streaming PDF: %v", err)
		return
	}
}
