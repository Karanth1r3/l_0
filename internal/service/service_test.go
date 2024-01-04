package service_test

// sorry for this one, may be refactored later with test cases

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Karanth1r3/l_0/internal/config"
	"github.com/Karanth1r3/l_0/internal/service"
	"github.com/Karanth1r3/l_0/internal/storage"
	"github.com/Karanth1r3/l_0/internal/utils"
	"github.com/google/uuid"
)

func TestService(t *testing.T) {
	t.Skip("For dev purpose only")

	type testR struct {
		orderUID     string
		body         []byte
		err          error
		expectedBody []byte
	}

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

func TestIndex(t *testing.T) {

	t.Skip("dev")
	req, err := http.NewRequest("GET", "/index", nil)
	if err != nil {
		t.Fatal(err)
	}

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
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.HandleIndex)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func Test404(t *testing.T) {

	t.Skip("DEV")

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

	orderUID := "wtf"

	req, err := http.NewRequest("GET", "/?id="+orderUID, nil)
	if err != nil {
		t.Fatal(err)
	}

	storage := storage.NewStorage(dbConn)
	storage.DropTable()

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
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.HandleGet)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

}

func TestGet(t *testing.T) {
	//	t.Skip("for dev purpose")

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
	orderUID := uuid.New().String()

	storage := storage.NewStorage(dbConn)
	storage.DropTable()

	err = storage.Write(orderUID, data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/?id="+orderUID, nil)
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
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.HandleGet)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	fmt.Println(rr.Body.String())
	// Check the response body is what we expect.
	if !bytes.Equal(data, []byte(rr.Body.String())) {
		t.Errorf("unexpected body: got %v want %v", rr.Body.String(), data)
	}
}

func TestPost(t *testing.T) {
	t.Skip("for dev purpose")

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
	orderUID := uuid.New().String()

	storage := storage.NewStorage(dbConn)
	storage.DropTable()

	err = storage.Write(orderUID, data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/?id="+orderUID, nil)
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
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.HandleGet)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}

}

func TestEmptyID(t *testing.T) {
	t.Skip("for dev purpose")

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
	orderUID := ""

	storage := storage.NewStorage(dbConn)
	storage.DropTable()

	err = storage.Write(orderUID, data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/?id="+orderUID, nil)
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
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.HandleGet)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

}
