package syncer

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// Service is a struct to wrap a running syncer
type Service struct {
	SourceSyncer      Syncer
	DestinationSyncer Syncer
	Access            Access
}

// NewService  init the service
func NewService(ss Syncer, ds Syncer, a Access) *Service {
	return &Service{
		SourceSyncer:      ss,
		DestinationSyncer: ds,
		Access:            a,
	}
}

// Run the service
func (s *Service) Run(ctx context.Context) error {
	log.Debugf("service.Run(ctx): ctx=%+v, s=%+v", ctx, s)
	return nil
}
