package syncer

import (
	"context"

	pgx "github.com/jackc/pgx/v4"
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

	// TODO: create a module to handle the pg actions and implement getTableColumns
	destinationColumns, err := s.getTableColumns(ctx, sourceConn, destinationConn)
	if err != nil {
		log.Errorf("service.Run(): s.getTableColumns() error=%w", err)
		return err
	}

	log.Debugf("service.Run(): destinationColumns: %+v", destinationColumns)

	// TODO: implements the clear_destination_table
	if s.Access.SyncMode == FullSync {
		err = s.truncateTable(ctx, destinationConn, s.Access.DestinationSchema, s.Access.DestinationTable)
		if err != nil {
			log.Errorf("service.Run(): s.truncateTable()  error=%w", err)
			return err
		}

		// TODO: move that to a function
		rows, err := sourceConn.Conn.Query(ctx, s.Access.SourceQuery)
		if err != nil {
			log.Errorf("service.Run(): sourceConn.Conn.Query(ctx, s.Access.SourceQuery) error=%w", err)
			return err
		}

		destinationIdentifier := pgx.Identifier{s.Access.DestinationSchema, s.Access.DestinationTable}

		copyCount, err := destinationConn.Conn.CopyFrom(ctx, destinationIdentifier, destinationColumns, rows)
		if err != nil {
			log.Errorf("service.Run(): sourceConn.Conn.CopyFrom() error=%w", err)
			return err
		}
		if rows.Err() != nil {
			log.Errorf("service.Run(): destinationConn.Conn.Query rows.Err()=%w", rows.Err())
			return rows.Err()
		}

		rows.Close()

		log.Debugf("service.Run(): copyCount: %+v", copyCount)
	}

	log.Debugf("service.Run(): sourceConn=%+v; destinationConn=%+v", sourceConn, destinationConn)

	return nil
}
