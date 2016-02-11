package main

// Some application constants.

const (
	// FIXME(alexyer): remove RESTAURANTS after refactoring.
	RESTAURANTS = "restaurants" // Restaurants database bucket
	QUEUES      = RESTAURANTS

	ORDERLIST = "orders"  // Orders list database bucket
	HISTORY   = "history" // History database bucket

	V1 = "v1" // Current API version

	// Command set
	CREATE_Q_CMD     = "create-q"
	DELETE_Q_CMD     = "delete-restaurant"
	CREATE_ORDER_CMD = "create-order"
	FINISH_ORDER_CMD = "finish-order"
	LIST_CMD         = "list"
	HISTORY_CMD      = "history"

	SUCCESS_RESPONSE = "success"
	ERROR_RESPONSE   = "error"
)
