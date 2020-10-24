package syncer

import (
	"context"
	"fmt"

	"github.com/cryp-com-br/pg-syncer/repository"
	pgx "github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
)

// getTableColumns
func (s *Service) getTableColumns(ctx context.Context, sourceConn *repository.PostgresConn, destinationConn *repository.PostgresConn, table string) (columns []string, err error) {
	rows, err := destinationConn.Conn.Query(ctx, "select column_name from information_schema.columns where table_name = $1", table)
	if err != nil {
		log.Errorf("service.Run(): destinationConn.Conn.Query error=%w", err)
		return
	}

	var c string
	for rows.Next() {
		rows.Scan(&c)
		columns = append(columns, c)
	}

	if rows.Err() != nil {
		log.Errorf("service.Run(): destinationConn.Conn.Query rows.Err()=%w", rows.Err())
		err = rows.Err()
		return
	}

	rows.Close()

	return
}

// truncateTable
func (s *Service) truncateTable(ctx context.Context, conn *repository.PostgresConn, schema, table string) (err error) {
	query := fmt.Sprintf("truncate table %s.%s", schema, table)
	_, err = conn.Conn.Exec(ctx, query)
	return
}

// copyFromSelect
func (s *Service) copyFromSelect(ctx context.Context, sourceConn *repository.PostgresConn, destinationConn *repository.PostgresConn) (err error) {
	// get columns
	destinationColumns, err := s.getTableColumns(ctx, sourceConn, destinationConn, s.Access.DestinationTable)
	if err != nil {
		log.Errorf("service.copyFromSelect(): s.getTableColumns() error=%w", err)
		return
	}

	// get the data to copy then
	rows, err := sourceConn.Conn.Query(ctx, s.Access.SourceQuery)
	if err != nil {
		log.Errorf("service.copyFromSelect(): sourceConn.Conn.Query(ctx, s.Access.SourceQuery) error=%w", err)
		return
	}

	destinationIdentifier := pgx.Identifier{s.Access.DestinationSchema, s.Access.DestinationTable}

	copyCount, err := destinationConn.Conn.CopyFrom(ctx, destinationIdentifier, destinationColumns, rows)
	if err != nil {
		log.Errorf("service.copyFromSelect(): sourceConn.Conn.CopyFrom() error=%w", err)
		return
	}
	if rows.Err() != nil {
		log.Errorf("service.copyFromSelect(): destinationConn.Conn.Query rows.Err()=%w", rows.Err())
		err = rows.Err()
		return
	}

	rows.Close()
	log.Debugf("service.copyFromSelect(): copyCount: %+v", copyCount)

	return
}
