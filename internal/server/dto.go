package server

import (
	"commuteboard/internal/domain"
	"time"
)

type CommuteResponse struct {
	OriginID        string    `json:"origin_id"`
	OriginName      string    `json:"origin_name"`
	DestinationID   string    `json:"destination_id"`
	DestinationName string    `json:"destination_name"`
	DurationMinutes int       `json:"duration_minutes"`
	DistanceKM      float64   `json:"distance_km"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewCommuteResponse(origin domain.Location, destination domain.Location, commute domain.Commute) CommuteResponse {
	return CommuteResponse{
		OriginID:        origin.ID,
		OriginName:      origin.Name,
		DestinationID:   destination.ID,
		DestinationName: destination.Name,
		DurationMinutes: int(commute.Duration.Minutes()),
		DistanceKM:      float64(commute.DistanceMeters) / 1000,
		UpdatedAt:       commute.RecordedAt,
	}
}
