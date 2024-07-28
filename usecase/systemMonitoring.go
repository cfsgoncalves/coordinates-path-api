package usecase

import (
	repositoryInterface "meight/repository/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SystemMonitoring struct {
	Cache repositoryInterface.Cache
}

func NewSystemMonitoring(cache repositoryInterface.Cache) *SystemMonitoring {
	return &SystemMonitoring{Cache: cache}
}

func (s *SystemMonitoring) Ping(context *gin.Context) {
	context.Status(http.StatusOK)
}

func (s *SystemMonitoring) Health(context *gin.Context) {
}
