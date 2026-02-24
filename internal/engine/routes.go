package engine

import (
	"bytes"
	"commuteboard/internal/domain"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	// "time"
)

var TRAVEL_MODE = "DRIVE"
var ROUTING_PREFERENCE = "TRAFFIC_AWARE_OPTIMAL"
var GOOGLE_ENDPOINT = "https://routes.googleapis.com/distanceMatrix/v2:computeRouteMatrix"

type matrixRequest struct {
	Origins           []matrixOrigin      `json:"origins"`
	Destinations      []matrixDestination `json:"destinations"`
	TravelMode        string              `json:"travelMode"`
	RoutingPreference string              `json:"routingPreference"`
	// DepartureTime     string              `json:"departureTime"`
}

type matrixOrigin struct {
	Waypoint matrixWaypoint `json:"waypoint"`
}

type matrixDestination struct {
	Waypoint matrixWaypoint `json:"waypoint"`
}

type matrixWaypoint struct {
	Location matrixLocation `json:"location"`
}

type matrixLocation struct {
	LatLng matrixLatLng `json:"latLng"`
}

type matrixLatLng struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func toMatrixOrigin(loc *domain.Location) (matrixOrigin, error) {
	lat, err := strconv.ParseFloat(loc.Latitude, 64)
	if err != nil {
		return matrixOrigin{}, err
	}
	lng, err := strconv.ParseFloat(loc.Longitude, 64)
	if err != nil {
		return matrixOrigin{}, err
	}

	origin := matrixOrigin{
		Waypoint: matrixWaypoint{
			Location: matrixLocation{
				LatLng: matrixLatLng{
					Latitude: lat, Longitude: lng,
				},
			},
		},
	}
	return origin, nil
}

func toMatrixDestination(loc *domain.Location) (matrixDestination, error) {
	lat, err := strconv.ParseFloat(loc.Latitude, 64)
	if err != nil {
		return matrixDestination{}, err
	}
	lng, err := strconv.ParseFloat(loc.Longitude, 64)
	if err != nil {
		return matrixDestination{}, err
	}

	destination := matrixDestination{
		Waypoint: matrixWaypoint{
			Location: matrixLocation{
				LatLng: matrixLatLng{
					Latitude: lat, Longitude: lng,
				},
			},
		},
	}
	return destination, nil
}

func (e *RouteEngine) computeRouteMatrix(
	ctx context.Context,
	destinations []*domain.Location,
) ([]domain.Route, error) {
	// Transform data into request body for API request
	origin, err := toMatrixOrigin(&e.Home)
	if err != nil {
		return nil, err
	}

	var matrixDestinations []matrixDestination
	for _, loc := range destinations {
		dest, err := toMatrixDestination(loc)
		if err != nil {
			return nil, err
		}
		matrixDestinations = append(matrixDestinations, dest)
	}

	reqBody := matrixRequest{
		Origins:           []matrixOrigin{origin},
		Destinations:      matrixDestinations,
		TravelMode:        TRAVEL_MODE,
		RoutingPreference: ROUTING_PREFERENCE,
		// DepartureTime:     time.Now().UTC().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, jsonData, "", "  "); err != nil {
		return nil, err
	}

	log.Printf("Matrix request:\n%s\n", prettyJSON.String())

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		GOOGLE_ENDPOINT,
		bytes.NewReader(jsonData),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", e.apiKey)
	req.Header.Set("X-Goog-FieldMask", "originIndex,destinationIndex,duration,distanceMeters")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Google response status: %s\n", resp.Status)
	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("Raw response:\n%s\n", string(bodyBytes))

	return []domain.Route{}, nil
}
