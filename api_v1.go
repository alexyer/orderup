package main

import (
	"encoding/json"
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

type listResponse struct {
	Response string   `json:"response"`
	Orders   *[]Order `json:"orders"`
}

func (o *Orderup) writeAPIResponse(w http.ResponseWriter, response []byte) {
	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}

// Encode and return error response to user.
func (o *Orderup) writeAPIErrorResponse(w http.ResponseWriter, apiErr error) {
	apiResponse := &APIErrorResponse{
		Response: ERROR_RESPONSE,
		Errors:   []string{apiErr.Error()},
	}

	response, err := json.Marshal(apiResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(response)
}

// list command.
func (o *Orderup) listAPIHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	req := &listRequest{}

	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orders, cmdErr := o.getPendingOrderList([]byte(req.Name))

	if cmdErr != nil {
		o.writeAPIErrorResponse(w, cmdErr)
		return
	}

	apiResponse := &listResponse{
		Response: "success",
		Orders:   orders,
	}

	response, err := json.Marshal(apiResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	o.writeAPIResponse(w, response)
}

// history command.
func (o *Orderup) historyAPIHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	req := &listRequest{}

	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orders, cmdErr := o.getHistoryList([]byte(req.Name))

	if cmdErr != nil {
		o.writeAPIErrorResponse(w, cmdErr)
		return
	}

	apiResponse := &listResponse{
		Response: "success",
		Orders:   orders,
	}

	response, err := json.Marshal(apiResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	o.writeAPIResponse(w, response)
}
