package prober

import (
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/pkg/prober/edge"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	appEdge *edge.Edge
	logger  *logrus.Entry
}

func newHandlers(
	appEdge *edge.Edge,
	logger *logrus.Entry,
) *Handlers {
	return &Handlers{
		appEdge: appEdge,
		logger:  logger,
	}
}

func (h *Handlers) Health(w http.ResponseWriter, req *http.Request) {
	healthy := h.appEdge.CheckHealth()
	if healthy != edge.StatusOk {
		h.logger.Error("Application is not healthy")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) Ready(w http.ResponseWriter, req *http.Request) {
	ready := h.appEdge.CheckReadiness()
	if ready != edge.StatusOk {
		h.logger.Error("Application is not ready")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
