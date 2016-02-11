package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (o *Orderup) getAPIv1() *API {
	return &API{
		Routes: []Route{
			Route{
				Path:        API_PREFIX + "/v1/queues/orders/list",
				HandlerFunc: o.listAPIHandler,
				Methods:     []string{"GET"},
			},
		},
	}
}

type listRequest struct {
	Name string `json:name`
}

func (o *Orderup) listAPIHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	req := &listRequest{}

	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println(req)
}
