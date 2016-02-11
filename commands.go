package main

// All bot command handlers located here.

import (
	"fmt"
	"strconv"
	"strings"
)

// create-q command.
// create-q [queue name]
func (o *Orderup) createQueueCmd(cmd *Cmd) (string, bool, *OrderupError) {
	switch {
	case len(cmd.Args) == 0:
		return "", true, NewOrderupError("Queue name is not given.", ARG_ERR)
	case len(cmd.Args) != 1:
		return "", true, NewOrderupError("Spaces are not allowed in queue name.", ARG_ERR)
	}

	name := cmd.Args[0]

	if err := o.createQueue([]byte(name)); err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	return fmt.Sprintf("%s queue created.", name), true, nil
}

// delete-q command.
// delete-q [queue name]
func (o *Orderup) deleteQueueCmd(cmd *Cmd) (string, bool, *OrderupError) {
	switch {
	case len(cmd.Args) == 0:
		return "", true, NewOrderupError("Queue name is not given.", ARG_ERR)
	case len(cmd.Args) != 1:
		return "", true, NewOrderupError("Spaces are not allowed in queue name.", ARG_ERR)
	}

	name := cmd.Args[0]

	if err := o.deleteQueue([]byte(name)); err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	return fmt.Sprintf("Queue: %s deleted.", name), true, nil
}

// create-order command.
// create-order [queue name] [@username] [order]
func (o *Orderup) createOrderCmd(cmd *Cmd) (string, bool, *OrderupError) {
	var (
		username   string
		name       string
		order      string
		orderCount int
		newOrder   *Order
		err        error
	)
	switch {
	case len(cmd.Args) < 3:
		return "", true, NewOrderupError("Wrong arguments", ARG_ERR)
	}

	username = cmd.Args[1]
	if username[0] != '@' {
		return "", true, NewOrderupError("Missing username", ARG_ERR)
	}

	name = cmd.Args[0]
	order = strings.Join(cmd.Args[2:], " ")

	newOrder, orderCount, err = o.createOrder([]byte(name), username, order)

	if err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	return fmt.Sprintf("%s order %d for %s %s - order %s. There are %d orders ahead of you.",
		name, newOrder.Id, newOrder.Username, order, order, orderCount), true, nil
}

// list command
// list [queue name]
func (o *Orderup) listCmd(cmd *Cmd) (string, bool, *OrderupError) {
	if len(cmd.Args) != 1 {
		return "", true, NewOrderupError("Wrong arguments", ARG_ERR)
	}

	queueName := cmd.Args[0]

	ordersList, err := o.getPendingOrderList([]byte(queueName))
	if err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	// Format orders list properly
	result := fmt.Sprintf("%s: pending orders:\n", queueName)
	for _, order := range *ordersList {
		result += order.String() + "\n"
	}

	return result, true, nil
}

// history command
// history [queue name]
func (o *Orderup) historyCmd(cmd *Cmd) (string, bool, *OrderupError) {
	if len(cmd.Args) != 1 {
		return "", true, NewOrderupError("Wrong arguments", ARG_ERR)
	}

	queueName := cmd.Args[0]

	history, err := o.getHistoryList([]byte(queueName))
	if err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	// Format orders list properly
	result := fmt.Sprintf("%s: history:\n", queueName)
	for _, order := range *history {
		result += order.String() + "\n"
	}

	return result, true, nil
}

// finish-order command
// finish-order [queue name] [order id]
func (o *Orderup) finishOrderCmd(cmd *Cmd) (string, bool, *OrderupError) {
	var (
		order *Order
	)

	if len(cmd.Args) != 2 {
		return "", true, NewOrderupError("Wrong arguments", ARG_ERR)
	}

	name := cmd.Args[0]
	orderId, err := strconv.Atoi(cmd.Args[1])
	if err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	order, err = o.finishOrder([]byte(name), orderId)

	if err != nil {
		return "", true, NewOrderupError(err.Error(), CMD_ERR)
	}

	return fmt.Sprintf("%s your order is finished. %s: Order: %d. %s",
		order.Username, name, order.Id, order.Order), true, nil
}

// help command.
func (o *Orderup) helpCmd(cmd *Cmd) (string, bool, *OrderupError) {
	return `Available commands:
				/orderup create-q [name] -- Create a list of order numbers for queue <name>.
				/orderup delete-q [name] -- Delete queue <name> and all orders in that queue.
				/orderup create-order [queue name] [@username] [order] -- Create a new order.
				/orderup finish-order [queue name]  [order id] -- Finish order.
				/orderup history [queue name] -- Show history for queue name.
				/orderup list [queue name] -- Get the list of orders for queue name.`, true, nil
}

// Helper function. Return error message with help contents.
func (o *Orderup) errorMessage(msg string) string {
	help, _, _ := o.helpCmd(nil)
	return fmt.Sprintf("%s\n%s", msg, help)
}
