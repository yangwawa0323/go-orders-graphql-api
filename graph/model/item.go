package model

// Item struct
type Item struct {
	ID          int    `json:"id"`
	ProductCode string `json:"productCode"`
	ProductName string `json:"productName"`
	Quantity    int    `json:"quantity"`
	OrderID		uint   `json:"-"`
}
