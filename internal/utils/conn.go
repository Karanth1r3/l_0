package utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Karanth1r3/l_0/internal/config"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

// Use params from config to make connstring & connect
func ConnectDB(cfg config.DB) (*sql.DB, error) {
	// connect to db with name "postgres"
	dbConnStr := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Name, cfg.Username, cfg.Password,
	)
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("connect to database failed: %w", err)
	}
	// set sessions limits
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(0)
	// try to ping db
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database is not available: %w", err)
	}

	// if everything's ok
	return db, nil
}

// Similar connection with config for nats streaming using it's api
func ConnectNats(cfg *config.NATS) (*nats.Conn, error) {
	conn, err := nats.Connect(
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		nats.Name("Service"),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to nats failed %w", err)
	}

	return conn, nil
}

func ConnectStan(nc *nats.Conn, cfg *config.STAN) (stan.Conn, error) {

	conn, err := stan.Connect(cfg.ClusterID, cfg.ClientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		return nil, fmt.Errorf("connection to nats streaming server failed: %w", err)
	}

	return conn, nil
}
