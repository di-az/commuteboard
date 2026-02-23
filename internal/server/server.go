package server

import (
	"commuteboard/internal/domain"
	"commuteboard/internal/store"
	"encoding/json"
	"log"
	"net/http"
)

const PORT = ":3333"

type HttpServer struct {
	store *store.RouteStore
}

func NewHttpServer(store *store.RouteStore) *HttpServer {
	return &HttpServer{store: store}
}

func (s *HttpServer) Run() {
	http.HandleFunc("/routes", s.GetRoutes)
	http.HandleFunc("/routes/{id}", s.GetRouteByID)

	http.ListenAndServe(PORT, nil)
}

func (s *HttpServer) GetRoutes(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting routes\n")
	var responseRoutes []domain.RouteResponse
	for _, route := range s.store.GetAll() {
		response := route.NewRouteResponse()
		responseRoutes = append(responseRoutes, response)
	}
	writeJSON(w, http.StatusOK, responseRoutes)
}

func (s *HttpServer) GetRouteByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("Getting route for id: %v\n", id)
	route, err := s.store.GetByID(id)
	routeResp := route.NewRouteResponse()
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, routeResp)
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// if data == nil {
	// 	return errors.New("empty response body")
	// }

	return json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, err error) {
	// log.Println(err.Error())
	_ = writeJSON(w, http.StatusInternalServerError, err.Error())
}
