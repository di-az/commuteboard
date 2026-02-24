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
	Home       domain.Location
	Locations  []*domain.Location
	Store      *store.RouteStore
	UpdateRate time.Duration
	TickRate   time.Duration
	running    atomic.Bool
	lastTick   atomic.Value
	client     *http.Client
	apiKey     string
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
	return &RouteEngine{
		Home:       home,
		Locations:  locations,
		Store:      store,
		UpdateRate: updateRate,
		TickRate:   tickRate,
		client:     &http.Client{Timeout: 5 * time.Second},
		apiKey:     apiKey,
	}
}

func (e *RouteEngine) checkLocations(ctx context.Context) {
	log.Printf("engine tick at %s", time.Now())
	now := time.Now()
	e.lastTick.Store(now)

	var toUpdate []*domain.Location

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

		toUpdate = append(toUpdate, location)
	}

	if len(toUpdate) == 0 {
		return
	}

	// route, err := getRoute(e.Home, *location)
	routes, err := e.computeRouteMatrix(ctx, toUpdate)
	if err != nil {
		log.Printf("error getting matrix: %v\n", err)
		return
	}

	for i, route := range routes {
		toUpdate[i].Schedule.LastUpdated = now
		e.Store.Set(route)
		log.Printf("Route updated: %s -> %s (%d min)",
			e.Home.Name,
			route.Finish.Name,
			route.Minutes,
		)
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

// TODO: Temporary Function
func getRoute(start, finish domain.Location) (domain.Route, error) {
	route := domain.Route{
		Start:        start,
		Finish:       finish,
		Minutes:      18,
		TrafficLevel: "Low",
		Timestamp:    time.Now(),
	}
	return route, nil
}
