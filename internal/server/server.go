package server

import (
	"commuteboard/internal/store"
	"fmt"
	"net/http"
)

const PORT = ":3333"

type HttpServer struct {
	store *store.RouteStore
}

func NewHttpServer(store *store.RouteStore) *HttpServer {
	return &HttpServer{store: store}
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func (s *HttpServer) GetRoutes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "getting routes\n")
	routes := s.store.GetAll()
	route := s.store.GetAll()[len(routes)-1]
	// for _, route := range s.store.GetAll() {
	fmt.Fprintf(w, "%v - %v: %v\n", route.Start.Name, route.Finish.Name, route.Minutes)
}

func (s *HttpServer) Run() {
	http.HandleFunc("/routes", s.GetRoutes)

	http.ListenAndServe(PORT, nil)
}
