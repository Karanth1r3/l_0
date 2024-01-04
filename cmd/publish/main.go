package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Karanth1r3/l_0/internal/config"
	"github.com/Karanth1r3/l_0/internal/model"
	"github.com/Karanth1r3/l_0/internal/utils"
	"github.com/google/uuid"
)

// made for adding records to storage

func main() {
	cfg, err := config.Parse("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	natsConn, err := utils.ConnectNats(&cfg.NATS)
	if err != nil {
		log.Fatal(err)
	}
	defer natsConn.Close()

	cfg.STAN.ClientID += "publisher" // this one has to differ from listener
	stanConn, err := utils.ConnectStan(natsConn, &cfg.STAN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := stanConn.Close(); err != nil {
			log.Println(err)
		}
	}()
	for i := 0; i < 5; i++ {

		orderUID := uuid.New().String()
		msg, err := json.Marshal(&model.Record{
			OrderUID: orderUID,
		})
		if err != nil {
			log.Fatal(err)
		}

		err = stanConn.Publish(
			cfg.Service.QueueName,
			msg,
		)
		if err != nil {
			log.Fatal(err)
		}
		//for debugging
		fmt.Println(orderUID)
		time.Sleep(time.Second)
	}
}
