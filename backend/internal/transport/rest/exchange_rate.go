package rest

import (
	"encoding/json"
	"net/http"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/service"
)

type createExchangeRateRequest struct {
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	Rate          string `json:"rate"`
	EffectiveDate string `json:"effective_date"`
	Source        string `json:"source"`
}

func createExchangeRateHandler(svc *service.ExchangeRateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createExchangeRateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		rate, err := svc.CreateRate(req.BaseCurrency, req.QuoteCurrency, req.Rate, req.EffectiveDate, req.Source)
		if err != nil {
			writeError(w, err)
			return
		}

		writeCreated(w, rate)
	}
}

func getExchangeRateHandler(svc *service.ExchangeRateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baseCurrency := r.PathValue("base")
		quoteCurrency := r.PathValue("quote")
		date := r.PathValue("date")

		rate, err := svc.GetRate(baseCurrency, quoteCurrency, date)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, rate)
	}
}

type listExchangeRatesHandlerWrapper struct {
	svc *service.ExchangeRateService
}

func listExchangeRatesHandler(svc *service.ExchangeRateService) http.HandlerFunc {
	return listExchangeRatesHandlerWrapper{svc: svc}.handle
}

func (h listExchangeRatesHandlerWrapper) handle(w http.ResponseWriter, r *http.Request) {
	baseCurrency := r.URL.Query().Get("base_currency")
	quoteCurrency := r.URL.Query().Get("quote_currency")

	pg := model.Pagination{
		Page:     parseInt(r.URL.Query().Get("page"), 1),
		PageSize: parseInt(r.URL.Query().Get("per_page"), parseInt(r.URL.Query().Get("page_size"), 200)),
	}

	result, err := h.svc.ListRates(baseCurrency, quoteCurrency, pg)
	if err != nil {
		writeError(w, err)
		return
	}

	// Frontend expects a flat array for exchange rates
	items := result.Items
	if items == nil {
		items = []model.ExchangeRate{}
	}
	writeOK(w, items)
}

func deleteExchangeRateHandler(svc *service.ExchangeRateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseInt64PathParam(r, "id")
		if err != nil {
			writeError(w, err)
			return
		}

		err = svc.DeleteRate(id)
		if err != nil {
			writeError(w, err)
			return
		}

		writeStatusNoContent(w)
	}
}

type fetchExchangeRatesRequest struct {
	BaseCurrency   string   `json:"base_currency"`
	QuoteCurrencies []string `json:"quote_currencies"`
}

func fetchExchangeRatesHandler(svc *service.ExchangeRateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req fetchExchangeRatesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, model.NewValidationError("invalid request body"))
			return
		}

		err := svc.FetchRatesFromAPI(req.BaseCurrency, req.QuoteCurrencies)
		if err != nil {
			writeError(w, err)
			return
		}

		writeStatusNoContent(w)
	}
}
