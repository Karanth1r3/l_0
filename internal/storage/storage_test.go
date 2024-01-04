package storage_test

import (
	"bytes"
	"testing"

	"github.com/Karanth1r3/l_0/internal/config"
	"github.com/Karanth1r3/l_0/internal/storage"
	"github.com/Karanth1r3/l_0/internal/utils"
)

func TestStorage(t *testing.T) {
	t.Skip("FOR DEV PURPOSE ONLY")

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

	storage := storage.NewStorage(dbConn)
	storage.DropTable()

	data := []byte(`{"a": "b"}`)
	orderUID := "trash"

	err = storage.Write(orderUID, data)
	if err != nil {
		t.Fatal(err)
	}
	result, err := storage.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	// if after table truncate & adding test record len is not 1, smthng went wrong
	if len(result) != 1 {
		t.Fatal("unexpected behaviour")
	}

	// if data is not equal to the one was supposed to be writed, smthng went wrong
	received := result[0]
	if !bytes.Equal(data, received.Data) {
		t.Fatal("unexpected behaviour")
	}
}
