package internal

import (
	"context"
	"github.com/tomwright/grace"
	"log"
	"time"
)

// NewCleanupRunner returns a runner that can be used cleanup old diagrams.
func NewCleanupRunner(generator Generator) grace.Runner {
	return &cleanupService{
		generator:   generator,
		runEvery:    time.Minute * 5,
		cleanupLast: time.Hour,
	}
}

// cleanupService is a runner that is used cleanup old diagrams.
type cleanupService struct {
	generator   Generator
	runEvery    time.Duration
	cleanupLast time.Duration
}

// Run starts the cleanup process.
func (s *cleanupService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := s.generator.CleanUp(s.cleanupLast); err != nil {
			log.Printf("error when cleaning up: %s", err.Error())
		}

		select {
		case <-time.After(s.runEvery):
			continue
		case <-ctx.Done():
			return nil
		}
	}
}
