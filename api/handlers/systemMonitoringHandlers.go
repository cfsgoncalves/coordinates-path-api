package api

import (
	repository "meight/repository/interfaces"
	"meight/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SystemMonitoringAPI struct {
	SystemMonitoring usecase.SystemMonitoring
}

func NewSystemMonitoringAPI(cache repository.Cache, database repository.Database, messageQueue repository.MessageQueue) *SystemMonitoringAPI {
	return &SystemMonitoringAPI{SystemMonitoring: *usecase.NewSystemMonitoring(cache, database, messageQueue)}
}

func (s *SystemMonitoringAPI) Ping(c *gin.Context) {
	if s.SystemMonitoring.Ping() {
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusInternalServerError)
}

func (s *SystemMonitoringAPI) Health(c *gin.Context) {
	if s.SystemMonitoring.Health() {
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusInternalServerError)
}
