package handler

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

// Pinger interface.
type Pinger interface {
	Ping(ctx context.Context) error
}

type ping struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

// Healthz handler.
type Healthz struct {
	*chi.Mux
	pingers map[string]Pinger
}

// Show handle GET /
func (h Healthz) Show(w http.ResponseWriter, r *http.Request) {
	var (
		wg     sync.WaitGroup
		status = 200
		pings  = make([]ping, len(h.pingers))
	)

	wg.Add(len(h.pingers))

	i := 0
	for service, pinger := range h.pingers {
		go func(i int, service string, pinger Pinger) {
			defer wg.Done()

			pings[i].Service = service
			if err := pinger.Ping(r.Context()); err != nil {
				logger.Error("ping error", zap.Error(err))

				status = 503
				pings[i].Status = err.Error()
			} else {
				pings[i].Status = "UP"
			}
		}(i, service, pinger)
		i++
	}
	wg.Wait()

	render(w, pings, status)
}

// Add a pinger.
func (h *Healthz) Add(name string, ping Pinger) {
	h.pingers[name] = ping
}

// NewHealthz handler.
func NewHealthz() Healthz {
	h := Healthz{
		Mux:     chi.NewMux(),
		pingers: make(map[string]Pinger),
	}

	h.Get("/", h.Show)

	return h
}
