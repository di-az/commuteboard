package main

import (
	"fmt"
	"log"
	"time"
)

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

type TimeRange struct {
	Start time.Duration
	End   time.Duration
}

type Schedule struct {
	Days map[time.Weekday][]TimeRange
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
		Minutes:      10,
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

	nowMinutes := time.Duration(t.Hour())*time.Hour +
		time.Duration(t.Minute())*time.Minute

	for _, timeRange := range dayRanges {
		if nowMinutes >= timeRange.Start && nowMinutes <= timeRange.End {
			log.Printf("Need to run now for: %v %v-%v", t.Weekday(), timeRange.Start, timeRange.End)
			return true
		}
	}

	return false
}

func main() {
	locations := []Location{work}
	t := time.Now()
	timeNow := time.Date(2026, 2, 21, 10, 30, 0, 0, t.Location())

	for _, location := range locations {
		location.Schedule.ShouldRunNow(timeNow)
		route, err := getRoute(house, work)
		if err != nil {
			log.Fatalf("Error calculating route: %v", err.Error())
		}
		fmt.Printf("Go from: %v to %v\n", house.Name, location.Name)
		fmt.Printf("Route will take: %v\n", route.Minutes)
	}
}
