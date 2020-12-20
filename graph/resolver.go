package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import "github.com/jinzhu/gorm"

// Resolver struct
type Resolver struct {
	DB *gorm.DB
}
