package transport

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/amiosamu/gofemart/internal/domain"
)

// @Summary OrderUploading
// @Description Загружает номер заказа в систему.
// @Security ApiKeyAuth
// @Tags orders
// @ID add order ID
// @Accept json
// @Param input body string true "order ID"
// @Success 202 "Status Accepted"
// @Failure 200 "Status OK"
// @Failure 400 "Bad Request"
// @Failure 401 "Status Unauthorized"
// @Failure 409 "Conflict"
// @Failure 422 "Status Unprocessable Entity"
// @Failure 500 "Internal Server Error"
// @Router /api/user/orders [post]
func (s *APIServer) OrderUploading(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logError("orderUploading", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.orders.AddOrderID(r.Context(), string(data)); err != nil {
		switch {
		case errors.Is(err, domain.ErrAlreadyUploadedByThisUser):
			logError("orderUploading", err)
			w.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, domain.ErrAlreadyUploadedByAnotherUser):
			logError("orderUploading", err)
			w.WriteHeader(http.StatusConflict)
			return
		case errors.Is(err, domain.ErrIncorrectOrder):
			logError("orderUploading", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		default:
			logError("orderUploading", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

// @Summary GetAllOrders
// @Description Выводит отсортированный по дате список заказов пользователя.
// @Security ApiKeyAuth
// @Tags orders
// @ID get all orders
// @Produce json
// @Success 200 {object} []domain.Order
// @Failure 204 "Status No Content"
// @Failure 401 "Status Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /api/user/orders [get]
func (s *APIServer) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := s.orders.GetAllOrders(r.Context())
	if err != nil {
		if errors.Is(err, domain.ErrNoData) {
			logError("getAllOrders", err)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		logError("getAllOrders", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ordersJSON, err := json.Marshal(orders)
	if err != nil {
		logError("getAllOrders", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ordersJSON)
}
