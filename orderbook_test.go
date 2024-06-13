package main

import (
	"fmt"
	"testing"
)

func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true, 12)
	buyOrderB := NewOrder(true, 3)
	buyOrderC := NewOrder(true, 7)
	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	fmt.Println(l)
	l.DeleteOrder(buyOrderB)
	fmt.Println(l)

	// askOrder := NewOrder(false, 5)
	// l.AddOrder(askOrder)
}

func TestOrderBook(t *testing.T) {
	ob := NewOrderBook()

	buyOrderA := NewOrder(true, 10)
	buyOrderB := NewOrder(true, 2000)
	ob.PlaceOrder(18_000, buyOrderA)
	ob.PlaceOrder(19_000, buyOrderB)

	fmt.Printf("%+v", ob)

	// We need to sort the orders by time
	// We need to sort the limit by price

}
