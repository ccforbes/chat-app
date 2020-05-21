package main

import (
	"log"
	"net/http"
	"os"

	"github.com/UW-Info-441-Winter-Quarter-2020/homework-cforbes1/servers/summary/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":5100"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
	log.Printf("server is listening at %s...", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
