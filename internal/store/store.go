package store

import (
	"commuteboard/internal/domain"
	"errors"
	"log"
	"sync"
)

type RouteStore struct {
	mu       sync.RWMutex
	commutes map[string]domain.Commute
}

func NewRouteStore() *RouteStore {
	return &RouteStore{commutes: make(map[string]domain.Commute)}
}

func (s *RouteStore) Set(commute domain.Commute) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := commute.DestinationID
	s.commutes[key] = commute
}

func (s *RouteStore) GetAll() []domain.Commute {
	s.mu.RLock()
	defer s.mu.RUnlock()

	log.Printf("GETTING ALL ROUTES: %v\n", s.commutes)
	result := make([]domain.Commute, 0, len(s.commutes))
	for _, c := range s.commutes {
		result = append(result, c)
	}
	return result
}

func (s *RouteStore) GetByID(id string) (domain.Commute, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	commute, ok := s.commutes[id]
	if !ok {
		return domain.Commute{}, errors.New("route not found")
	}

	return commute, nil
}

func (s *RouteStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.commutes, id)
}
