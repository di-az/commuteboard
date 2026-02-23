package server

import (
	"commuteboard/internal/domain"
	"time"
)

type RouteResponse struct {
	DestinationID   string    `json:"destination_id"`
	DestinationName string    `json:"destination_name"`
	Minutes         int       `json:"minutes"`
	TrafficLevel    string    `json:"traffic_level"`
	Timestamp       time.Time `json:"timestamp"`
}

func NewRouteResponse(r domain.Route) RouteResponse {
	return RouteResponse{
		DestinationID:   r.Finish.ID,
		DestinationName: r.Finish.Name,
		Minutes:         r.Minutes,
		TrafficLevel:    r.TrafficLevel,
		Timestamp:       r.Timestamp,
	}
}
