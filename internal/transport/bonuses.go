package transport

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/amiosamu/gofemart/internal/domain"
)

// @Summary Balance
// @Description Выводит сумму баллов лояльности и использованных за весь период регистрации баллов пользователя.
// @Security ApiKeyAuth
// @Tags balance
// @ID balance
// @Produce json
// @Success 200 {object} domain.BalanceOutput
// @Failure 401 "Status Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /api/user/balance [get]
func (s *APIServer) Balance(w http.ResponseWriter, r *http.Request) {
	balance, err := s.withdraw.Balance(r.Context())
	if err != nil {
		logError("balance", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		logError("balance", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(balanceJSON)
}

// @Summary Withdraw
// @Description Реализует списание бонусов пользователя в учет суммы нового заказа.
// @Security ApiKeyAuth
// @Tags withdraw
// @ID withdraw
// @Accept json
// @Param input body domain.Withdraw true "Запрос параметров списания"
// @Success 200 "OK"
// @Failure 401 "Status Unauthorized"
// @Failure 402 "Status Payment Required"
// @Failure 422 "Status Unprocessable Entity"
// @Failure 500 "Internal Server Error"
// @Router /api/user/balance/withdraw [post]
func (s *APIServer) Withdraw(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logError("withdraw", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var withdraw domain.Withdraw
	if err := json.Unmarshal(data, &withdraw); err != nil {
		logError("withdraw", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.withdraw.Withdraw(r.Context(), withdraw); err != nil {
		if errors.Is(err, domain.ErrAlreadyUploadedByThisUser) {
			logError("withdraw", err)
			w.WriteHeader(http.StatusOK)
			return
		} else if errors.Is(err, domain.ErrAlreadyUploadedByAnotherUser) {
			logError("withdraw", err)
			w.WriteHeader(http.StatusConflict)
			return
		} else if errors.Is(err, domain.ErrIncorrectOrder) {
			logError("withdraw", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		} else if errors.Is(err, domain.ErrNoBonuses) {
			logError("withdraw", err)
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}
		logError("withdraw", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// @Summary Withdrawals
// @Description Выводит отсортированный по дате список списаний бонусов пользователя.
// @Tags withdraw
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} []domain.Withdraw
// @Failure 204 "Status No Content"
// @Failure 401 "Status Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /api/user/withdrawals [get]
func (s *APIServer) Withdrawals(w http.ResponseWriter, r *http.Request) {
	withdrawals, err := s.withdraw.Withdrawals(r.Context())
	if err != nil {
		if errors.Is(err, domain.ErrNoWithdraws) {
			logError("withdrawals", err)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		logError("withdrawals", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdrawalsJSON, err := json.Marshal(withdrawals)
	if err != nil {
		logError("withdrawals", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(withdrawalsJSON)
}
