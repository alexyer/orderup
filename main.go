package main

import (
	"io"
	"net/http"
	"strings"
)

// Command struct.
type Cmd struct {
	Name string   // Command name
	Args []string // List of command arguments
}

// Handle requests to orderup bot.
func orderup(w http.ResponseWriter, r *http.Request) {
	// Parse request
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get command text from the request and split arguments.
	cmd := parseCmd(r.PostForm["text"][0])

	// Execute command
	response := execCmd(cmd)

	io.WriteString(w, response)
}

func parseCmd(cmd string) *Cmd {
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

func execCmd(cmd *Cmd) string {
	switch cmd.Name {
	case "create-restaurant":
		return createRestaurant(cmd)
	default:
		return `Available commands:
					/orderup create-restaurant [name] -- Create a list of order numbers for restaurant name.`
	}
}

func main() {
	http.HandleFunc("/", orderup)
	http.ListenAndServe(":4242", nil)
}
