package common

import "time"

type StockUpdate struct {
	Source string
	Sku    string
	Now    time.Time
	Qty    int
}

type ReindexRequest struct {
	Source string
	Sku    string
	Client string
}

type ExportRequest struct {
	Aggregate string
	Sku       string
	Client    string
}
