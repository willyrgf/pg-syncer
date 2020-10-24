package syncer

import (
	"context"
	"fmt"

	"github.com/cryp-com-br/pg-syncer/repository"
	log "github.com/sirupsen/logrus"
)

// getTableColumns
func (s *Service) getTableColumns(ctx context.Context, sourceConn *repository.PostgresConn, destinationConn *repository.PostgresConn) (columns []string, err error) {
	rows, err := destinationConn.Conn.Query(ctx, "select column_name from information_schema.columns where table_name = $1", s.Access.DestinationTable)
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
func (s *Service) truncateTable(ctx context.Context, conn *repository.PostgresConn, table, schema string) (err error) {
	query := fmt.Sprintf("truncate table %s.%s", schema, table)
	_, err = conn.Conn.Exec(ctx, query)
	return
}
