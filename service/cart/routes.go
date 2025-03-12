package cart

import (
	"fmt"
	"net/http"

	"github.com/Wirlhawk/go-ecommerce/service/auth"
	"github.com/Wirlhawk/go-ecommerce/types"
	"github.com/Wirlhawk/go-ecommerce/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	store        types.OrderStore
	productStore types.ProductStore
	userStore    types.UserStore
}

func NewHandler(store types.OrderStore, productStore types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, productStore: productStore, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/cart/checkout", auth.WithJWTAuth(h.handleCheckout, h.userStore)).Methods(http.MethodPost)
}

func (h *Handler) handleCheckout(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())

	// get cart payload
	var cart types.CartCheckoutPayload

	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid message field %v", error))
		return
	}

	productIDs, err := getCardItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ps, err := h.productStore.GetProductsByIDs(productIDs)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := h.createOrder(ps, cart.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, map[string]any{
		"total_price": totalPrice,
		"order_id":    orderID,
	})
}
