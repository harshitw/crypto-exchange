package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	Size      float64 // size of order
	Limit     *Limit  // pointer to limit
	Timestamp int64   // use unix nano
	Bid       bool    // is bid or ask
}

type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

// Limit is a group of orders at a certain price level
type Limit struct {
	Price       float64
	Orders      Orders
	TotalVolume float64
}

type Limits []*Limit

type ByBestAsk struct{ Limits }

func (a ByBestAsk) Len() int           { return len(a.Limits) }
func (a ByBestAsk) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsk) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestBid struct{ Limits }

func (b ByBestBid) Len() int           { return len(b.Limits) }
func (b ByBestBid) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestBid) Less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{}, // we cannot preallocate exact space for orders
	}
}

func (l *Limit) String() string {
	return fmt.Sprintf("[price: %2f | volume: %2f]", l.Price, l.TotalVolume)
}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1] // shift total size by 1
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}
	o.Limit = nil
	l.TotalVolume -= o.Size

	// To Do : resort the whole orders
	sort.Sort(l.Orders)
}

type Orderbook struct {
	Asks []*Limit
	Bids []*Limit

	AskLimits map[float64]*Limit // prices and points to limit
	BidLimits map[float64]*Limit
}

func NewOrderBook() *Orderbook {
	return &Orderbook{
		Asks:      []*Limit{},
		Bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

// Placing an order will always return a slice of matches but they could be empty
func (ob *Orderbook) PlaceOrder(price float64, o *Order) []Match {
	// 1. try to match the order
	// matching logic

	// 2. add the rest of the orders to the books
	if o.Size > 0.0 {
		ob.add(price, o)
	}
	return []Match{}
}

// IF BUY or SELL 10 BTC at 15K
func (ob *Orderbook) add(price float64, o *Order) {
	// check if volume sitting at price
	var limit *Limit

	// Bids
	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		limit.AddOrder(o)
		if o.Bid {
			ob.Bids = append(ob.Bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.Asks = append(ob.Asks, limit)
			ob.AskLimits[price] = limit
		}
	}

}
