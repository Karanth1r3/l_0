package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Karanth1r3/l_0/internal/cache"
	"github.com/Karanth1r3/l_0/internal/config"
	"github.com/Karanth1r3/l_0/internal/model"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
)

const queryParamId = "id"

type (
	Service struct {
		cache    srvCache
		storage  srvStorage
		htmlPage []byte
		sub      stan.Subscription
	}

	srvCache interface {
		Write(orderUID string, data []byte)
		Read(orderUID string) ([]byte, error)
	}

	srvStorage interface {
		Write(orderUID string, data []byte) error
		ReadAll() ([]model.StorageRecord, error)
	}
)

// Constructor
func NewService(storage srvStorage, conn stan.Conn, cfg config.Service, htmlPage []byte) (*Service, error) {
	s := &Service{
		cache:    cache.NewCache(),
		storage:  storage,
		htmlPage: htmlPage,
	}

	// Process subscriber options
	startOpt := stan.StartAt(pb.StartPosition_NewOnly)
	if cfg.StartSeq != 0 {
		startOpt = stan.StartAtSequence(cfg.StartSeq)
	} else if cfg.DeliverLast {
		startOpt = stan.StartWithLastReceived()
	} else if cfg.DeliverAll && !cfg.NewOnly {
		startOpt = stan.DeliverAllAvailable()
	} else if cfg.StartDelta != "" {
		ago, err := time.ParseDuration(cfg.StartDelta)
		if err != nil {
			return nil, fmt.Errorf("failed to parse start delta: %w", err)
		}
		startOpt = stan.StartAtTimeDelta(ago)
	}

	sub, err := conn.QueueSubscribe(
		cfg.QueueName,
		cfg.QueueGroup,
		s.HandleMessage,
		startOpt,
		stan.DurableName("Service"))
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe queue: %w", err)
	}
	s.sub = sub
	return s, s.UpdateCache()
}

// Handle message from queue
func (s *Service) HandleMessage(msg *stan.Msg) {

	var r model.Record

	err := json.Unmarshal(msg.Data, &r)
	if err != nil {
		log.Println("failed to unmarshal record", err)
		return
	}

	if !r.IsValid() {
		log.Println("record is invalid (empty orderUID)")
		return
	}

	err = s.Write(r.OrderUID, msg.Data)

}

// Read data from cache
func (s *Service) Read(orderUID string) ([]byte, error) {
	return s.cache.Read(orderUID)
}

// Write data to db & cache
func (s *Service) Write(orderUID string, data []byte) error {
	if err := s.storage.Write(orderUID, data); err != nil {
		return fmt.Errorf("failed to write to db: %w", err)
	}
	s.cache.Write(orderUID, data)

	return nil
}

// Restore data from database to cache
func (s *Service) UpdateCache() error {

	//Read all json values from db
	all, err := s.storage.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read from db: %w", err)
	}

	// for any values from db
	for _, r := range all {
		s.cache.Write(r.OrderUID, r.Data)
	}

	return nil
}

// Main Page handler
func (s *Service) HandleIndex(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write(s.htmlPage)
}

// HandleGet handles get request by id in query
func (s *Service) HandleGet(w http.ResponseWriter, r *http.Request) {

	// if req is not GET, show according status
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	// extract id from query out of request body/
	id := r.URL.Query().Get(queryParamId)
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}

	record, err := s.cache.Read(id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(http.StatusText(http.StatusNotFound)))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	_, _ = w.Write(record)
}
