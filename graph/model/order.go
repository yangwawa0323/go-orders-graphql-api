package model

// Order struct
type Order struct {
	ID           int     `json:"id" gorm:"primary_key"`
	CustomerName string  `json:"customerName"`
	OrderAmount  float64 `json:"orderAmount"`
	Items        []*Item `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:OrderID" json:"items"`
}
