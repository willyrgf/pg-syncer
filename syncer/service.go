package syncer

import (
	"context"
	"fmt"

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
	rows, err := destinationConn.Conn.Query(ctx, "select column_name from information_schema.columns where table_name = $1", s.Access.DestinationTable)
	if err != nil {
		log.Errorf("service.Run(): destinationConn.Conn.Query error=%w", err)
		return err
	}

	var destinationColumns []string
	for rows.Next() {
		var c string
		rows.Scan(&c)
		destinationColumns = append(destinationColumns, c)
	}

	if rows.Err() != nil {
		log.Errorf("service.Run(): destinationConn.Conn.Query rows.Err()=%w", rows.Err())
		return rows.Err()
	}

	rows.Close()

	log.Debugf("service.Run(): destinationColumns: %+v", destinationColumns)
	// TODO: implements the clear_destination_table
	if s.Access.CleanDestinationTable {
		query := fmt.Sprintf("truncate table %s.%s", s.Access.DestinationSchema, s.Access.DestinationTable)
		_, err = destinationConn.Conn.Exec(ctx, query)
		if err != nil {
			log.Errorf("service.Run(): destinationConn.Conn.Query('truncate table') error=%w", err)
			return err
		}
	}

	// TODO: move that to a function
	rows, err = sourceConn.Conn.Query(ctx, s.Access.SourceQuery)
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

	log.Debugf("service.Run(): sourceConn=%+v; destinationConn=%+v", sourceConn, destinationConn)

	return nil
}
