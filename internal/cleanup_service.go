package internal

import (
	"log"
	"time"
)

// NewCleanupService returns a service that can be used cleanup old diagrams.
func NewCleanupService(generator Generator) *cleanupService {
	return &cleanupService{
		generator: generator,
		stopCh:    make(chan struct{}),
	}
}

// cleanupService is a service that can be used cleanup old diagrams.
type cleanupService struct {
	generator Generator
	stopCh    chan struct{}
}

// Start starts the cleanup service.
func (s *cleanupService) Start() error {
	for {
		if err := s.generator.CleanUp(time.Hour); err != nil {
			log.Printf("error when cleaning up: %s", err.Error())
		}

		select {
		case <-time.After(time.Minute * 5):
			continue
		case <-s.stopCh:
			return nil
		}
	}
}

// Stop stops the cleanup service.
func (s *cleanupService) Stop() {
	close(s.stopCh)
}
