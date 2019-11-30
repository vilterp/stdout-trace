package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	mux *http.ServeMux

	socketsMu sync.Mutex
	sockets   []*websocket.Conn

	messagesMu sync.Mutex
	messages   []string
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
	http.FileServer(http.Dir("trace-ui/build")).ServeHTTP(w, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) serveWS(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	s.addSocket(conn)

	s.catchUp(conn)

	select {} // block
}

func (s *Server) catchUp(conn *websocket.Conn) {
	s.messagesMu.Lock()
	defer s.messagesMu.Unlock()

	for _, msg := range s.messages {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Println("failed to catch up socket: ", err)
		}
	}
}

func (s *Server) addSocket(socket *websocket.Conn) {
	s.socketsMu.Lock()
	defer s.socketsMu.Unlock()

	s.sockets = append(s.sockets, socket)

}

// TODO: consider using a file instead of keeping this all in memory
func (s *Server) appendMessage(line string) {
	s.messagesMu.Lock()
	defer s.messagesMu.Unlock()

	s.messages = append(s.messages, line)
}

func (s *Server) processStdin() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		s.appendMessage(line)
		s.pushToSockets(line)
		fmt.Println(line)
	}
	fmt.Println("eof")
}

func (s *Server) pushToSockets(line string) {
	// even parse it?
	for idx, socket := range s.sockets {
		if err := socket.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			log.Println("failed to write to socket at idx", idx)
		}
	}
}
