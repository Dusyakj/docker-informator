package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"project8/internal/handler"
	"project8/pkg/serv"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: server <address>")
		os.Exit(1)
	}

	addr := os.Args[1]

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/image-download-size", handler.GetImageInformation)

	srv, err := serv.NewServer(addr, mux)
	if err != nil {
		log.Fatalf("Ошибка создания сервера: %v", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatalf("Ошибка при работе сервера: %v", err)
	}

}
