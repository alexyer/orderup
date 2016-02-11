package main

// Some application constants.

const (
	QUEUES       = "queues"
	ORDERLIST    = "orders"        // Orders list database bucket
	HISTORY      = "history"       // History database bucket
	ORDERCOUNTER = "order_counter" // Indicates the last order id

	V1 = "v1" // Current API version

	// Command set
	CREATE_Q_CMD     = "create-q"
	DELETE_Q_CMD     = "delete-q"
	CREATE_ORDER_CMD = "create-order"
	FINISH_ORDER_CMD = "finish-order"
	LIST_CMD         = "list"
	HISTORY_CMD      = "history"

	SUCCESS_RESPONSE = "success"
	ERROR_RESPONSE   = "error"
)
