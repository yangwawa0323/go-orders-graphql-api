package model

// Order struct
type Order struct {
	ID           int     `json:"id" gorm:"primary_key"`
	CustomerName string  `json:"customerName"`
	OrderAmount  float64 `json:"orderAmount"`
	Items        []*Item `gorm:"foreignKey:OrderID" json:"items"`
}
