package cart

import (
	"fmt"

	"github.com/Wirlhawk/go-ecommerce/types"
)

func getCardItemsIDs(items []types.CartItem) ([]int, error) {
	productIDs := make([]int, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for the product %d", item.ProductID)
		}

		productIDs[i] = item.ProductID
	}

	return productIDs, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItem, userID int) (int, float64, error) {
	productMap := make(map[int]types.Product)
	for _, product := range ps {
		productMap[product.ID] = product
	}

	// check if products in stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, nil
	}

	// calculate total price
	totalPrice := calculateTotalPrice(items, productMap)

	// reduce qty of product in db
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity

		_ = h.productStore.UpdateProduct(product)
	}

	// create the order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address",
	})
	if err != nil {
		return 0, 0, err
	}

	for _, item := range items {
		_ = h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	// create order_items
	return orderID, totalPrice, nil
}

func checkIfCartIsInStock(cartItems []types.CartItem, products map[int]types.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the store, please refresh your cart", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not available in the quantity requested", item.ProductID)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []types.CartItem, products map[int]types.Product) float64 {
	var total float64
	for _, item := range cartItems {
		product := products[item.ProductID]
		total += product.Price * float64(item.Quantity)
	}

	return total
}
