package usecase

import (
	repositoryInterface "meight/repository/interfaces"
)

type SystemMonitoring struct {
	Cache        repositoryInterface.Cache
	Database     repositoryInterface.Database
	MessageQueue repositoryInterface.MessageQueue
}

func NewSystemMonitoring(cache repositoryInterface.Cache, db repositoryInterface.Database, ms repositoryInterface.MessageQueue) *SystemMonitoring {
	return &SystemMonitoring{Cache: cache, Database: db, MessageQueue: ms}
}

func (s *SystemMonitoring) Ping() bool {
	return true
}

// Need to implemente Health function
func (s *SystemMonitoring) Health() bool {
	return s.Cache.Ping() && s.Database.Ping() && s.MessageQueue.Ping()
}
