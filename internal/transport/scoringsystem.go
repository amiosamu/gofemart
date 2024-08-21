package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/amiosamu/gofemart/internal/domain"
)

func (s *APIServer) ScoringSystem() {
	orderID, err := s.scoringsystem.GetOrderStatus(context.Background())
	if err != nil {
		logError("scoringSystem", err)
		return
	}

	for _, id := range orderID {

		addr := fmt.Sprintf("%s/api/orders/%s", s.config.ScoringSystemPort, id)
		resp, err := http.Get(addr)
		if err != nil {
			logError("scoringSystem", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				logError("scoringSystem", err)
				return
			}

			var orderScoring domain.ScoringSystem
			if err := json.Unmarshal(data, &orderScoring); err != nil {
				logError("scoringSystem", err)
				return
			}

			if err := s.scoringsystem.UpdateOrder(context.Background(), orderScoring); err != nil {
				logError("scoringSystem", err)
				return
			}
		}
	}

}
