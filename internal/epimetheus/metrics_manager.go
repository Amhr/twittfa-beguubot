package epimetheus

import "github.com/cafebazaar/epimetheus"

type MetricsManager struct {
	Ep       *epimetheus.Epimetheus
	Requests *epimetheus.TimerWithCounter
	Counter  *epimetheus.Counter
}

func NewMetricsManager(ep *epimetheus.Epimetheus) *MetricsManager {
	return &MetricsManager{
		Ep: ep,
		Requests: ep.NewTimerWithCounter("Requests", []string{
			"type",
			"route",
			"status",
		}),
		Counter: ep.NewCounter("Counter", []string{"type", "path", "status"}),
	}
}

func (mm *MetricsManager) CacheNotExists(route string) {
	mm.Counter.Inc("cache", route, "not-found")
}

func (mm *MetricsManager) CacheExists(route string) {
	mm.Counter.Inc("cache", route, "ok")
}
