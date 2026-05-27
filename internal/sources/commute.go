package sources

import (
	"fmt"
	"signalboard/internal/content"
	"signalboard/internal/engine"
)

type CommuteSource struct {
	Engine *engine.RouteEngine
}

func NewCommuteSource(engine *engine.RouteEngine) *CommuteSource {
	return &CommuteSource{
		Engine: engine,
	}
}

func (c *CommuteSource) Content() []content.Content {
	routes := c.Engine.GetRoutes()

	var items []content.Content

	for _, route := range routes {
		if route.DurationSeconds == nil {
			continue
		}

		minutes := int(route.DurationSeconds.Minutes())
		label := fmt.Sprintf("%s → %s", route.Origin.Name, route.Destination.Name)

		items = append(items, content.Metric{
			Label: label,
			Value: fmt.Sprintf("%d", minutes),
			Unit:  "min",
		})
	}

	return items
}
