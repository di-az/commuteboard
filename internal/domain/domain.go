package domain

import (
	"time"
)

type TimeRange struct {
	Start time.Duration
	End   time.Duration
}

type Schedule struct {
	Days        map[time.Weekday][]TimeRange
	LastUpdated time.Time
}

type Location struct {
	ID        string
	Name      string
	Latitude  string
	Longitude string
	Schedule  Schedule
}

type Route struct {
	Start        Location
	Finish       Location
	Minutes      int
	TrafficLevel string
	Timestamp    time.Time
}

type RouteResponse struct {
	DestinationID   string    `json:"destination_id"`
	DestinationName string    `json:"destination_name"`
	Minutes         int       `json:"minutes"`
	TrafficLevel    string    `json:"traffic_level"`
	Timestamp       time.Time `json:"timestamp"`
}

func (s Schedule) ShouldRunNow(t time.Time) bool {
	dayRanges, ok := s.Days[t.Weekday()]
	if !ok {
		return false
	}

	nowMinutes := time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute

	for _, timeRange := range dayRanges {
		if nowMinutes >= timeRange.Start && nowMinutes <= timeRange.End {
			// log.Printf("Need to run now for: %v %v-%v", t.Weekday(), timeRange.Start, timeRange.End)
			return true
		}
	}

	return false
}

func (r *Route) NewRouteResponse() RouteResponse {
	return RouteResponse{
		DestinationID:   r.Finish.ID,
		DestinationName: r.Finish.Name,
		Minutes:         r.Minutes,
		TrafficLevel:    r.TrafficLevel,
		Timestamp:       r.Timestamp,
	}
}
