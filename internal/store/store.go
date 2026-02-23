package store

import (
	"commuteboard/internal/domain"
	"errors"
	"sync"
)

type RouteStore struct {
	mu     sync.RWMutex
	routes map[string]domain.Route
}

func NewRouteStore() *RouteStore {
	return &RouteStore{routes: make(map[string]domain.Route)}
}

func (s *RouteStore) Set(route domain.Route) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := route.Finish.ID
	s.routes[key] = route
}

func (s *RouteStore) GetAll() []domain.Route {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.Route, 0, len(s.routes))
	for _, r := range s.routes {
		result = append(result, r)
	}
	return result
}

func (s *RouteStore) GetByID(id string) (domain.Route, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	route, ok := s.routes[id]
	if !ok {
		return domain.Route{}, errors.New("route not found")
	}

	return route, nil
}
