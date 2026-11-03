package rest

import (
	"encoding/json"
	"net/http"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/service"
)

func getBalanceSheetHandler(svc *service.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")
		depth := parseInt(r.URL.Query().Get("depth"), 0)

		var tagFilters []int64
		if tagsStr := r.URL.Query().Get("tag_filters"); tagsStr != "" {
			if err := json.Unmarshal([]byte(tagsStr), &tagFilters); err != nil {
				writeError(w, model.NewValidationError("invalid tag_filters: "+err.Error()))
				return
			}
		}

		report, err := svc.GetBalanceSheet(dateFrom, dateTo, depth, tagFilters)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, report)
	}
}

func getIncomeStatementHandler(svc *service.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")
		depth := parseInt(r.URL.Query().Get("depth"), 0)

		var tagFilters []int64
		if tagsStr := r.URL.Query().Get("tag_filters"); tagsStr != "" {
			if err := json.Unmarshal([]byte(tagsStr), &tagFilters); err != nil {
				writeError(w, model.NewValidationError("invalid tag_filters: "+err.Error()))
				return
			}
		}

		report, err := svc.GetIncomeStatement(dateFrom, dateTo, depth, tagFilters)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, report)
	}
}

func getNetWorthHandler(svc *service.ReportService) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        asOfDate := r.URL.Query().Get("as_of_date")
        if asOfDate == "" {
            writeError(w, model.NewValidationError("as_of_date is required"))
            return
        }

        netWorth, err := svc.GetNetWorth(asOfDate)
        if err != nil {
            writeError(w, err)
            return
        }

        writeOK(w, netWorth)
    }
}

func getCashFlowHandler(svc *service.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")

		cashFlow, err := svc.GetCashFlow(dateFrom, dateTo)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, map[string]string{
			"cash_flow": cashFlow.String(),
		})
	}
}

func getTopSpendingHandler(svc *service.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")
		limit := parseInt(r.URL.Query().Get("limit"), 10)

		categories, err := svc.GetTopSpendingCategories(dateFrom, dateTo, limit)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, categories)
	}
}

func getExpenseBreakdownHandler(svc *service.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")
		groupBy := r.URL.Query().Get("group_by")

		if groupBy == "" {
			groupBy = "account"
		}

		breakdown, err := svc.GetExpenseBreakdown(dateFrom, dateTo, groupBy)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, breakdown)
	}
}

func getTagSpendingSummaryHandler(svc *service.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagKey := r.URL.Query().Get("tag_key")
		dateFrom := r.URL.Query().Get("date_from")
		dateTo := r.URL.Query().Get("date_to")

		if tagKey == "" {
			writeError(w, model.NewValidationError("tag_key is required"))
			return
		}

		summary, err := svc.GetTagSpendingSummary(tagKey, dateFrom, dateTo)
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, summary)
	}
}

func getNetWorthHandlerWithTrend(svc *service.ReportService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		asOfDate := r.URL.Query().Get("as_of_date")
		if asOfDate == "" {
			writeError(w, model.NewValidationError("as_of_date is required"))
			return
		}

		comparisonDate := r.URL.Query().Get("comparison_date")

		var netWorth *service.NetWorthWithTrend
		var err error
		if comparisonDate != "" {
			netWorth, err = svc.GetNetWorth(asOfDate, comparisonDate)
		} else {
			netWorth, err = svc.GetNetWorth(asOfDate)
		}
		if err != nil {
			writeError(w, err)
			return
		}

		writeOK(w, netWorth)
	}
}
