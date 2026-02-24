package engine

import (
	"commuteboard/internal/domain"
	"commuteboard/internal/store"
	"context"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type RouteEngine struct {
	Home          domain.Location
	Locations     []*domain.Location
	locationIndex map[string]*domain.Location
	Store         *store.RouteStore
	UpdateRate    time.Duration
	TickRate      time.Duration
	running       atomic.Bool
	lastMeasured  map[string]time.Time
	lastTick      atomic.Value
	client        *http.Client
	apiKey        string
}

type Status struct {
	Running    bool   `json:"runnning"`
	TickRate   string `json:"tick_rate"`
	UpdateRate string `json:"update_rate"`
	Locations  int    `json:"locations"`
	LastTick   int    `json:"last_tick"`
}

func NewRouteEngine(
	home domain.Location,
	locations []*domain.Location,
	store *store.RouteStore,
	updateRate time.Duration,
	tickRate time.Duration,
	apiKey string,
) *RouteEngine {
	locationIndex := make(map[string]*domain.Location)
	locationIndex[home.ID] = &home

	for _, loc := range locations {
		locationIndex[loc.ID] = loc
	}

	return &RouteEngine{
		Home:          home,
		Locations:     locations,
		locationIndex: locationIndex,
		Store:         store,
		UpdateRate:    updateRate,
		TickRate:      tickRate,
		lastMeasured:  make(map[string]time.Time),
		client:        &http.Client{Timeout: 5 * time.Second},
		apiKey:        apiKey,
	}
}

func (e *RouteEngine) checkLocations(ctx context.Context) {
	log.Printf("engine tick at %s", time.Now())
	now := time.Now()
	e.lastTick.Store(now)

	var activeDestinations []*domain.Location

	for _, location := range e.Locations {
		// Skip if not in time range
		if !location.Schedule.ShouldRunNow(now) {
			e.Store.Delete(location.ID)
			continue
		}

		// Skip if recently updated
		if now.Sub(location.Schedule.LastUpdated) < e.UpdateRate {
			continue
		}

		activeDestinations = append(activeDestinations, location)
	}

	if len(activeDestinations) == 0 {
		return
	}

	commutes, err := e.computeRouteMatrix(ctx, activeDestinations)
	if err != nil {
		log.Printf("error computing matrix: %v\n", err)
		return
	}

	for _, comm := range commutes {
		// activeDestinations[i].Schedule.LastUpdated = now
		e.lastMeasured[comm.DestinationID] = now
		log.Printf("Setting routes: %v\n", comm)
		e.Store.Set(comm)
	}
}

func (e *RouteEngine) Run(ctx context.Context) {
	e.running.Store(true)
	defer e.running.Store(false)

	ticker := time.NewTicker(e.TickRate)
	defer ticker.Stop()

	log.Printf("Route engine started\n")
	e.checkLocations(ctx)

	for {
		select {
		case <-ticker.C:
			e.checkLocations(ctx)
		case <-ctx.Done():
			log.Println("Route engine shutting down")
			return
		}
	}
}

func (e *RouteEngine) Status() Status {
	return Status{
		Running:    e.running.Load(),
		TickRate:   e.TickRate.String(),
		UpdateRate: e.UpdateRate.String(),
		Locations:  len(e.Locations),
	}
}

func (e *RouteEngine) GetLocationByID(id string) *domain.Location {
	return e.locationIndex[id]
}
