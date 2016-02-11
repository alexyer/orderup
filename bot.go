package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/boltdb/bolt"
)

// Command struct.
type Cmd struct {
	Name string   // Command name
	Args []string // List of command arguments
}

type Orderup struct {
	db *bolt.DB
}

func NewOrderup(dbFile string) (*Orderup, error) {
	db, err := initDb(dbFile)
	if err != nil {
		return nil, err
	}

	return &Orderup{
		db: db,
	}, nil
}

// Serve web API.
func (o *Orderup) makeAPI(apiVersion string, mux *http.ServeMux) {
	switch apiVersion {
	case V1:
		for _, route := range o.getAPIv1().Routes {
			mux.HandleFunc(route.Path, route.HandlerFunc)
		}

	default:
		panic("Unknown API version.")
	}
}

// Serve Slack API.
func (o *Orderup) makeRequestHandler(mux *http.ServeMux) {
	mux.HandleFunc("/orderup", o.requestHandler)
}

// Open an initialize database.
func initDb(dbFile string) (*bolt.DB, error) {
	db, err := bolt.Open(dbFile, 0600, nil)

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(RESTAURANTS))
		return err
	})

	return db, err
}

// Handle requests to orderup bot.
func (o *Orderup) requestHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get command text from the request and split arguments.
	cmd := o.parseCmd(r.PostForm["text"][0])

	// Execute command
	if response, inChannel := o.execCmd(cmd); inChannel {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, fmt.Sprintf(`{"response_type":"in_channel","text":"%s"}`, response))
	} else {
		io.WriteString(w, response)
	}
}

// Parse command from the request string.
func (o *Orderup) parseCmd(cmd string) *Cmd {
	if cmdLst := strings.Split(cmd, " "); len(cmdLst) == 1 {
		return &Cmd{
			Name: cmdLst[0],
		}
	} else {
		return &Cmd{
			Name: cmdLst[0],
			Args: cmdLst[1:],
		}
	}
}

// Execute command.
func (o *Orderup) execCmd(cmd *Cmd) (string, bool) {
	switch cmd.Name {
	case CREATE_Q_CMD:
		return o.createRestaurant(cmd)
	case DELETE_Q_CMD:
		return o.deleteRestaurant(cmd)
	case CREATE_ORDER_CMD:
		return o.createOrder(cmd)
	case FINISH_ORDER_CMD:
		return o.finishOrder(cmd)
	case LIST_CMD:
		return o.list(cmd)
	case HISTORY_CMD:
		return o.history(cmd)
	default:
		return o.help(cmd)
	}
}

// Safely close db and shutdown.
func (o *Orderup) Shutdown() {
	o.db.Close()
	log.Print("Bye!")
}
