package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
)

var (
	host     string
	port     int
	db       string
	password string
)

func init() {
	flag.StringVar(&host, "host", "localhost", "host to listen")
	flag.IntVar(&port, "port", 5000, "port to listen")
	flag.StringVar(&db, "db", "orderup.db", "database file")
	flag.StringVar(&password, "passcode", "", "protection password")
}

func main() {
	flag.Parse()

	// Create new instance of bot.
	bot, err := NewOrderup(db, password)
	if err != nil {
		log.Fatal(err)
	}

	// Configure shutdown process.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for _ = range c {
			bot.Shutdown()
			os.Exit(0)
		}
	}()

	mux := mux.NewRouter()

	bot.makeAPI(V1, mux)        // Make API handlers in the mux
	bot.makeRequestHandler(mux) // Make Slack API endpoint

	log.Printf("Orderup started on %s:%d.\n", host, port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux))
}
