package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
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

type Server struct {
	mux       *http.ServeMux
	socketsMu sync.Mutex
	sockets   []Socket
}

type Socket interface {
	Write(line string) error
}

func NewServer() *Server {
	s := &Server{}
	mux := http.NewServeMux()

	mux.Handle("/ws", http.HandlerFunc(s.serveWS))
	mux.Handle("/", http.HandlerFunc(s.serveStatic))

	s.mux = mux
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.mux.ServeHTTP(w, req)
}

func (s *Server) serveStatic(w http.ResponseWriter, req *http.Request) {
	http.FileServer(http.Dir("build")).ServeHTTP(w, req)
}

func (s *Server) serveWS(w http.ResponseWriter, req *http.Request) {
	s.socketsMu.Lock()
	defer s.socketsMu.Unlock()

	s.sockets = append(s.sockets, newSocket(w))
}

// ???
type httpSocket struct {
	w http.ResponseWriter
}

func (s *httpSocket) Write(line string) error {
	_, err := s.w.Write([]byte(line))
	return err
}

func newSocket(w http.ResponseWriter) *httpSocket {
	return &httpSocket{
		w: w,
	}
}

func (s *Server) processStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		s.pushToSockets(line)
		fmt.Print(line)
	}
}

func (s *Server) pushToSockets(line string) {
	// even parse it?
	for idx, socket := range s.sockets {
		if err := socket.Write(line); err != nil {
			log.Println("failed to write to socket at idx", idx)
		}
	}
}
