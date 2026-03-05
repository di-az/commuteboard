package server

import (
	"commuteboard/internal/domain"
	"time"
)

type RouteResponse struct {
	ID              int        `json:"id"`
	Origin          string     `json:"origin"`
	Destination     string     `json:"destination"`
	DurationMinutes *int       `json:"duration_minutes"`
	DistanceKM      *float64   `json:"distance_km"`
	RecordedAt      *time.Time `json:"updated_at"`
}

func NewRouteResponse(route domain.Route) RouteResponse {
	r := RouteResponse{
		ID:          route.ID,
		Origin:      route.Origin.Name,
		Destination: route.Destination.Name,
	}

	if route.DurationSeconds != nil {
		min := int(route.DurationSeconds.Minutes())
		r.DurationMinutes = &min
	}
	if route.DistanceMeters != nil {
		km := float64(*route.DistanceMeters) / 1000
		r.DistanceKM = &km
	}
	if route.RecordedAt != nil {
		r.RecordedAt = route.RecordedAt
	}

	return r
}
