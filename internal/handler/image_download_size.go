package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"project8/internal/requester"
)

func GetImageInformation(w http.ResponseWriter, r *http.Request) {
	log.Printf("API Request: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	output, code, err := requester.GetInformation(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(output)
}
