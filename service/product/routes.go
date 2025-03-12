package product

import (
	"fmt"
	"net/http"

	"github.com/Wirlhawk/go-ecommerce/types"
	"github.com/Wirlhawk/go-ecommerce/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)

}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, ps)
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.ProductPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

    // validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid message field %v", error))
		return
	}

	err := h.store.CreateProduct(types.Product{
		Name:        payload.Name,
		Description: payload.Description,
		Image:       payload.Image,
		Price:       payload.Price,
		Quantity:    payload.Quantity,
	})
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"data":    payload,
		"message": "Product created successfully",
		"success": true,
	})
}
