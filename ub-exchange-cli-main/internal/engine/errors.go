package engine

import "fmt"

var (
	ErrOrderExists = fmt.Errorf("orderbook: order already exists")
)
