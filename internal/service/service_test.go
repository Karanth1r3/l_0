package service_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/Karanth1r3/l_0/internal/config"
	"github.com/Karanth1r3/l_0/internal/service"
	"github.com/Karanth1r3/l_0/internal/storage"
	"github.com/Karanth1r3/l_0/internal/utils"
)

func TestService(t *testing.T) {
	//	t.Skip("For dev purpose only")

	dbConn, err := utils.ConnectDB(config.DB{
		Host:     "localhost",
		Port:     5432,
		Name:     "level0",
		Username: "service",
		Password: "q1w2e3",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer dbConn.Close()

	data := []byte(`{"a": "b"}`)
	orderUID := "trash"

	storage := storage.NewStorage(dbConn)
	storage.DropTable()

	err = storage.Write(orderUID, data)
	if err != nil {
		t.Fatal(err)
	}
	// nats mock
	natsConn, err := utils.ConnectNats(&config.NATS{
		Host: "localhost",
		Port: 4222,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer natsConn.Close()

	stanConn, err := utils.ConnectStan(natsConn, &config.STAN{
		ClusterID: "Level0",
		ClientID:  "ServiceTest",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := stanConn.Close(); err != nil {
			log.Println(err)
		}
	}()

	service, err := service.NewService(storage, stanConn, config.Service{
		QueueName:   "Queue",
		QueueGroup:  "QueueGroup",
		StartSeq:    0,
		DeliverLast: true,
		DeliverAll:  false,
		NewOnly:     true,
		StartDelta:  "0",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	err = service.Write(orderUID, data)
	if err != nil {
		t.Fatal(err)
	}
	received, err := service.Read(orderUID)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, received) {
		t.Fatal("unexpected behaviour")
	}
}
