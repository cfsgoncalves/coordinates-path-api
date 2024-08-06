package usecase

import (
	repositoryInterface "meight/repository/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SystemMonitoring struct {
	Cache        repositoryInterface.Cache
	Database     repositoryInterface.Database
	MessageQueue repositoryInterface.MessageQueue
}

func NewSystemMonitoring(cache repositoryInterface.Cache, db repositoryInterface.Database, ms repositoryInterface.MessageQueue) *SystemMonitoring {
	return &SystemMonitoring{Cache: cache, Database: db, MessageQueue: ms}
}

func (s *SystemMonitoring) Ping(context *gin.Context) {
	context.Status(http.StatusOK)
}

// Need to implemente Health function
func (s *SystemMonitoring) Health(context *gin.Context) {
}
