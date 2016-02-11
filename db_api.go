package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
)

// DB calls.
// Available for Slack and REST APIs.

// Get orders for the queue.
// General function to get list of orders from the queue.
// bucket can be ORDERLIST or HISTORY.
func (o *Orderup) getOrderList(queueName []byte, bucket []byte) (*[]Order, error) {
	var orderList []Order

	err := o.db.View(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(QUEUES))

		queue := b.Bucket(queueName)
		if queue == nil {
			return NonExistentQueue(string(queueName))
		}

		orders := queue.Bucket(bucket)
		c := orders.Cursor()

		// Iterate over all orders, decode and store in the orders list
		for k, v := c.First(); k != nil; k, v = c.Next() {
			order := Order{}
			if err := json.Unmarshal(v, &order); err != nil {
				return err
			}

			orderList = append(orderList, order)
		}

		return
	})

	if err != nil {
		return nil, err
	}

	return &orderList, nil
}

// Get pending orders list for the queue.
func (o *Orderup) getPendingOrderList(queue []byte) (*[]Order, error) {
	return o.getOrderList(queue, []byte(ORDERLIST))
}

// Get history list for the queue.
func (o *Orderup) getHistoryList(queue []byte) (*[]Order, error) {
	return o.getOrderList(queue, []byte(HISTORY))

}

// Create new queue <name>.
func (o *Orderup) createQueue(name []byte) error {
	return o.db.Update(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(QUEUES))

		// Create new bucket for the new queue.
		r, err := b.CreateBucket(name)
		if err != nil {
			return errors.New("Queue already exists.")
		}

		// Create 2 subbucktes.
		// One for pending orders, another for finished orders.
		_, err = r.CreateBucket([]byte(ORDERLIST))
		_, err = r.CreateBucket([]byte(HISTORY))
		err = r.Put([]byte(ORDERCOUNTER), itob(0))

		return err
	})
}

// Delete queue <name>.
func (o *Orderup) deleteQueue(name []byte) error {
	return o.db.Update(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(QUEUES))

		// Create new bucket for the new queue.
		if err = b.DeleteBucket(name); err != nil {
			return errors.New("Restaurant does not exist.")
		}

		return err
	})
}

// Create <order> in the <queue> for <username>.
// Return order id and orders count.
func (o *Orderup) createOrder(queueName []byte, username, order string) (int, int, error) {
	var (
		id         uint64
		orderCount int
	)

	err := o.db.Update(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(QUEUES))

		queue := b.Bucket(queueName)
		if queue == nil {
			return NonExistentQueue(string(queueName))
		}

		orders := queue.Bucket([]byte(ORDERLIST))

		// Prepare order data
		id, _ = orders.NextSequence()
		orderCount = orders.Stats().KeyN

		// JSON serialize order
		buf, err := json.Marshal(&Order{
			Username: username,
			Order:    order,
			Id:       int(id),
		})

		if err != nil {
			return err
		}

		// Update order counter
		queue.Put([]byte(ORDERCOUNTER), itob(int(id)))

		// Store order into the database
		return orders.Put(itob(int(id)), buf)
	})

	return int(id), orderCount, err
}

// Finish order <orderId> in the <queue>.
func (o *Orderup) finishOrder(queueName []byte, orderId int) (*Order, error) {
	var (
		order     Order
		orderData []byte
	)

	err := o.db.Batch(func(tx *bolt.Tx) (err error) {
		// Get bucket with restaurants.
		b := tx.Bucket([]byte(QUEUES))

		queue := b.Bucket(queueName)
		if queue == nil {
			return errors.New(fmt.Sprintf("Restaurant %s does not exist", queue))
		}

		orders := queue.Bucket([]byte(ORDERLIST))
		history := queue.Bucket([]byte(HISTORY))
		orderCounter := int(btoi(queue.Get([]byte(ORDERCOUNTER))))

		if orderId > orderCounter {
			return errors.New("Too big order id. Order does not exist yet.")
		}

		orderData = orders.Get(itob(orderId))
		if orderData == nil {
			return errors.New("Order is already finished.")
		}

		// Delete order from the orders list
		if err := orders.Delete(itob(orderId)); err != nil {
			return err
		}

		// Put order in the history list
		return history.Put(itob(orderId), orderData)
	})

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(orderData, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
