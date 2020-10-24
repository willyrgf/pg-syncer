package syncer

import (
	"context"
	"fmt"

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

	sourceConn, err := s.SourceSyncer.Repo.GetConn()
	if err != nil {
		log.Errorf("service.Run(): s.SourceSyncer.Repo.GetConn() error=%w", err)
		return err
	}
	defer sourceConn.Conn.Release()

	destinationConn, err := s.DestinationSyncer.Repo.GetConn()
	if err != nil {
		log.Errorf("service.Run(): s.SourceSyncer.Repo.GetConn() error=%w", err)
		return err
	}
	defer destinationConn.Conn.Release()

	log.Debugf("service.Run(): sourceConn=%+v; destinationConn=%+v", sourceConn, destinationConn)

	switch s.Access.SyncMode {
	case FullSync:
		log.Debugf("sync_mode selected: %s", FullSync)

		err = s.truncateTable(ctx, destinationConn, s.Access.DestinationSchema, s.Access.DestinationTable)
		if err != nil {
			log.Errorf("service.Run(): s.truncateTable()  error=%w", err)
			return err
		}

		err = s.copyFromSelect(ctx, sourceConn, destinationConn)
		if err != nil {
			log.Errorf("service.Run(): s.copyFromSelect()  error=%w", err)
			return err
		}
	default:
		log.Errorf("the sync_mode configured is not implemented yet: %s", s.Access.SyncMode)
		return fmt.Errorf("the sync_mode configured is not implemented yet: %s", s.Access.SyncMode)

	}

	return nil
}
