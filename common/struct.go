package common

import "time"

type StockUpdate struct {
	Source string
	Sku string
	Now time.Time
	Qty int
}