package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("listening on http://%s/", addr)

	s := NewServer()
	go s.processStdin()
	log.Fatal(http.ListenAndServe(addr, s))
}
