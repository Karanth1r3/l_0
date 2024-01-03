package main

import (
	"log"
	"os"

	"github.com/Karanth1r3/l_0/internal/config"
	"github.com/Karanth1r3/l_0/internal/utils"
)

func main() {
	// trying to initialize config based on provided file
	cfg, err := config.Parse("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	// []byte data from page to feed to the server (i guess)
	htmlPage, err := os.ReadFile("index.html")
	if err != nil {
		log.Fatal(err)
	}
	// setting db connection was put to an internal module for convinience
	dbConn, err := utils.ConnectDB(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { // on main goroutine executed cast conn close check
		if err := dbConn.Close(); err != nil {
			log.Println(err)
		}
	}()
	// connect to nats streaming service using config params & defer it's closure on end
	natsConn, err := utils.ConnectNats(&cfg.NATS)
	if err != nil {
		log.Fatal(err)
	}
	defer natsConn.Close()

	stanConn, err := utils.ConnectStan(natsConn, &cfg.STAN)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := stanConn.Close(); err != nil {
			log.Println(err)
		}
	}()
}
