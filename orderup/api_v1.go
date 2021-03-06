package orderup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Get APIv1 layout.
func (o *Orderup) getAPIv1() *API {
	return &API{
		Routes: []Route{
			Route{
				Path:        API_PREFIX + "/v1/queues/orders/list",
				HandlerFunc: o.BasicAuth(o.listAPIHandler),
				Methods:     []string{"GET"},
			},
			Route{
				Path:        API_PREFIX + "/v1/queues/orders/history",
				HandlerFunc: o.BasicAuth(o.historyAPIHandler),
				Methods:     []string{"GET"},
			},
			Route{
				Path:        API_PREFIX + "/v1/queues",
				HandlerFunc: o.BasicAuth(o.createQueueAPIHandler),
				Methods:     []string{"POST"},
			},
			Route{
				Path:        API_PREFIX + "/v1/queues",
				HandlerFunc: o.BasicAuth(o.deleteQueueAPIHandler),
				Methods:     []string{"DELETE"},
			},
			Route{
				Path:        API_PREFIX + "/v1/queues/order",
				HandlerFunc: o.BasicAuth(o.createOrderAPIHandler),
				Methods:     []string{"POST"},
			},
			Route{
				Path:        API_PREFIX + "/v1/queues/orders/finish",
				HandlerFunc: o.BasicAuth(o.finishOrderAPIHandler),
				Methods:     []string{"PUT"},
			},
		},
	}
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

// Actual API call - createQueue, getHistoryList, etc.
// Decode request into map, executes API call and returns response.
type apiAction func(map[string]string) (interface{}, error)

func (o *Orderup) apiHandler(w http.ResponseWriter, r *http.Request, exec apiAction) {
	// Decode API request.
	decoder := json.NewDecoder(r.Body)

	req := make(map[string]string)

	if err := decoder.Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Execute API call
	if response, cmdErr := exec(req); cmdErr != nil {
		o.writeAPIErrorResponse(w, cmdErr)
		return
	} else {
		buf, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		o.writeAPIResponse(w, buf)
	}
}

// list command.
func (o *Orderup) listAPIHandler(w http.ResponseWriter, r *http.Request) {
	o.apiHandler(w, r, func(req map[string]string) (interface{}, error) {
		if req["name"] == "" {
			return nil, WrongArgsError()
		}

		orders, cmdErr := o.getPendingOrderList([]byte(req["name"]))
		if cmdErr != nil {
			return nil, cmdErr
		}

		return struct {
			Response string   `json:"response"`
			Orders   *[]Order `json:"orders"`
		}{
			Response: SUCCESS_RESPONSE,
			Orders:   orders,
		}, nil
	})
}

// history command.
func (o *Orderup) historyAPIHandler(w http.ResponseWriter, r *http.Request) {
	o.apiHandler(w, r, func(req map[string]string) (interface{}, error) {
		if req["name"] == "" {
			return nil, WrongArgsError()
		}

		orders, cmdErr := o.getHistoryList([]byte(req["name"]))
		if cmdErr != nil {
			return nil, cmdErr
		}

		return struct {
			Response string   `json:"response"`
			Orders   *[]Order `json:"orders"`
		}{
			Response: SUCCESS_RESPONSE,
			Orders:   orders,
		}, nil
	})
}

// create-q command.
func (o *Orderup) createQueueAPIHandler(w http.ResponseWriter, r *http.Request) {
	o.apiHandler(w, r, func(req map[string]string) (interface{}, error) {
		if req["name"] == "" {
			return nil, WrongArgsError()
		}

		cmdErr := o.createQueue([]byte(req["name"]))
		if cmdErr != nil {
			return nil, cmdErr
		}

		return struct {
			Response string
		}{
			Response: fmt.Sprintf("Queue %s created.", req["name"]),
		}, nil
	})
}

// delete-q command.
func (o *Orderup) deleteQueueAPIHandler(w http.ResponseWriter, r *http.Request) {
	o.apiHandler(w, r, func(req map[string]string) (interface{}, error) {
		if req["name"] == "" {
			return nil, WrongArgsError()
		}

		cmdErr := o.deleteQueue([]byte(req["name"]))
		if cmdErr != nil {
			return nil, cmdErr
		}

		return struct {
			Response string
		}{
			Response: fmt.Sprintf("Queue %s deleted.", req["name"]),
		}, nil
	})
}

// create-order command.
func (o *Orderup) createOrderAPIHandler(w http.ResponseWriter, r *http.Request) {
	o.apiHandler(w, r, func(req map[string]string) (interface{}, error) {
		if req["name"] == "" || req["user"] == "" || req["description"] == "" {
			return nil, WrongArgsError()
		}

		req["user"] = "@" + req["user"]

		order, orderCount, cmdErr := o.createOrder([]byte(req["name"]), req["user"], req["description"])
		if cmdErr != nil {
			return nil, cmdErr
		}

		return struct {
			Response    string `json:"response"`
			Order       *Order `json:"order"`
			OrdersAhead int    `json:"orders_ahead"`
		}{
			Response:    SUCCESS_RESPONSE,
			Order:       order,
			OrdersAhead: orderCount,
		}, nil
	})
}

// finish-order command.
func (o *Orderup) finishOrderAPIHandler(w http.ResponseWriter, r *http.Request) {
	o.apiHandler(w, r, func(req map[string]string) (interface{}, error) {
		if req["name"] == "" || req["id"] == "" {
			return nil, WrongArgsError()
		}

		orderId, err := strconv.Atoi(req["id"])
		if err != nil {
			return nil, err
		}

		order, err := o.finishOrder([]byte(req["name"]), orderId)
		if err != nil {
			return nil, err
		}

		return struct {
			Response string `json:"response"`
			Order    *Order `json:"order"`
		}{
			Response: SUCCESS_RESPONSE,
			Order:    order,
		}, nil
	})
}
