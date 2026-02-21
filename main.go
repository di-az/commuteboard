package main

import (
	"log"
	"time"
)

const UpdateRate = 1 * time.Minute
const tickRate = 10 * time.Second

var house = Location{
	Name:      "House",
	Latitude:  "20.745317326696103",
	Longitude: "-103.44431208289149",
	// Schedule:  Schedule{Times: []string{"08:00-10:00"}},
}

var work = Location{
	Name:      "Work",
	Latitude:  "20.688900217575455",
	Longitude: "-103.42880959994349",
	Schedule: Schedule{
		Days: map[time.Weekday][]TimeRange{
			time.Tuesday: {
				{Start: 8 * time.Hour, End: 10 * time.Hour},
			},
			time.Thursday: {
				{Start: 8 * time.Hour, End: 10 * time.Hour},
			},
			time.Saturday: {
				{Start: 1 * time.Hour, End: 23 * time.Hour},
			},
		},
	},
}

var piano = Location{
	Name:      "Piano",
	Latitude:  "20.688900217575455",
	Longitude: "-103.42880959994349",
	Schedule: Schedule{
		Days: map[time.Weekday][]TimeRange{
			time.Saturday: {
				{Start: 9 * time.Hour, End: 18 * time.Hour},
			},
		},
	},
}

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

func getRoute(start, finish Location) (Route, error) {
	route := Route{
		Start:        start,
		Finish:       finish,
		Minutes:      18,
		TrafficLevel: "Low",
		Timestamp:    time.Now(),
	}
	return route, nil
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

func main() {
	locations := []*Location{&work, &piano}
	ticker := time.NewTicker(tickRate)
	defer ticker.Stop()

	log.Printf("Route engine has started\n")

	// t := time.Now()
	// timeNow := time.Date(2026, 2, 21, 10, 30, 0, 0, t.Location())

	for range ticker.C {
		log.Printf("Ticking")
		now := time.Now()
		for _, location := range locations {
			// Skip if not in time range
			if !location.Schedule.ShouldRunNow(now) {
				continue
			}

			// Skip if recently updated
			if now.Sub(location.Schedule.LastUpdated) < UpdateRate {
				continue
			}

			route, err := getRoute(house, *location)
			if err != nil {
				log.Printf("Error calculating route for %s: %v", location.Name, err)
				continue
			}

			location.Schedule.LastUpdated = now
			log.Printf("Route updated: %s -> %s (%d min)",
				house.Name,
				location.Name,
				route.Minutes,
			)
		}
	}
}
