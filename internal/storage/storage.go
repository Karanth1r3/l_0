package storage

import (
	"database/sql"
	"fmt"

	"github.com/Karanth1r3/l_0/internal/model"
)

const tableName = "service_storage"

type Storage struct {
	db *sql.DB
}

// Constructor
func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// Try to write to db (storage) [k - orderUID, v - orderInfo in []byte (data is static)
func (s *Storage) Write(orderUID string, data []byte) error {
	query := fmt.Sprintf(`INSERT INTO %s (order_uid, value)
	VALUES ($S1, $S2)
	ON CONFLICT (order_uid) DO UPDATE SET
	value = EXCLUDED.value;`, tableName) // rewrite on collision

	_, err := s.db.Exec(query, orderUID, data)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

// Try to read all data from db
func (s *Storage) ReadAll() ([]model.StorageRecord, error) {
	query := fmt.Sprintf(`SELECT order_uid, value FROM %s ;`, tableName)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query: %w", err)
	}

	result := make([]model.StorageRecord, 0)
	// iterator for query results
	for rows.Next() {
		var (
			data     []byte
			orderUID string
		)

		err = rows.Scan(&orderUID, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result = append(result, model.StorageRecord{
			OrderUID: orderUID,
			Data:     data,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return result, nil
}

// Remove all records from table (For testing)
func (s *Storage) DropTable() error {
	query := fmt.Sprintf(`TRUNCATE TABLE %s;`, tableName)

	_, err := s.db.Exec(query)

	if err != nil {
		return fmt.Errorf("failed to execute truncate query: %w", err)
	}

	return nil
}
