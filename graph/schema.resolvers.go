package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/yangwawa0323/go-orders-graphql-api/graph/generated"
	"github.com/yangwawa0323/go-orders-graphql-api/graph/model"
)

func (r *mutationResolver) CreateOrder(ctx context.Context, input model.OrderInput) (*model.Order, error) {
	// panic(fmt.Errorf("not implemented"))
	order := model.Order{
		CustomerName: input.CustomerName,
		OrderAmount:  input.OrderAmount,
		Items:        mapItemsFromInput(input.Items),
	}

	err := r.DB.Create(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *mutationResolver) UpdateOrder(ctx context.Context, orderID int, input model.OrderInput) (*model.Order, error) {
	// panic(fmt.Errorf("not implemented"))
	updatedOrder := model.Order{
		ID:           orderID,
		CustomerName: input.CustomerName,
		OrderAmount:  input.OrderAmount,
		Items:        mapItemsFromInput(input.Items),
	}
	err := r.DB.Save(&updatedOrder).Error //.Save(&updatedOrder.Items).Error
	if err != nil {
		return nil, err
	}
	r.DB.Save(&updatedOrder.Items)
	return &updatedOrder, nil
}

func (r *mutationResolver) DeleteOrder(ctx context.Context, orderID int) (bool, error) {
	// panic(fmt.Errorf("not implemented"))
	r.DB.Where("order_id = ?", orderID).Delete(&model.Item{})
	r.DB.Where("id = ?", orderID).Delete(&model.Order{})
	return true, nil
}

func (r *queryResolver) Orders(ctx context.Context) ([]*model.Order, error) {
	// panic(fmt.Errorf("not implemented"))
	var orders []*model.Order
	err := r.DB.Preload("Items").Find(&orders).Error
	// err := r.DB.Set("gorm:auto_preload", true).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
type itemResolver struct{ *Resolver }
type orderResolver struct{ *Resolver }

func mapItemsFromInput(itemsInput []*model.ItemInput) []*model.Item {
	var items []*model.Item
	for _, itemInput := range itemsInput {
		if itemInput.ID == nil {
			items = append(items, &model.Item{
				ProductCode: itemInput.ProductCode,
				ProductName: itemInput.ProductName,
				Quantity:    itemInput.Quantity,
			})
		} else {
			items = append(items, &model.Item{
				ID:          *itemInput.ID,
				ProductCode: itemInput.ProductCode,
				ProductName: itemInput.ProductName,
				Quantity:    itemInput.Quantity,
			})
		}
	}
	return items
}
