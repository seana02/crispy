package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crispy-2/backend/internal/service"
	"github.com/crispy-2/backend/internal/testutil"
)

func setupTransactionRestTest(t *testing.T) (*service.AccountService, *service.TransactionService, *httptest.Server) {
	t.Helper()
	db := testutil.NewTestDB(t)
	acctSvc := service.NewAccountService(db)
	rateSvc := service.NewExchangeRateService(db)
	txSvc := service.NewTransactionService(db, rateSvc)
	txSvc.SetAccountService(acctSvc)
	mux := SetupRouter(db, acctSvc, txSvc, nil, nil, nil, nil, nil, nil, nil)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return acctSvc, txSvc, server
}

// Regression: update transaction must use PUT, not PATCH.
// Before fix, the frontend sent PATCH which returned 404 because only PUT is registered.
func TestUpdateTransaction_UsePUT(t *testing.T) {
	acctSvc, _, server := setupTransactionRestTest(t)

	// Create accounts
	assets, _ := acctSvc.CreateAccount("Assets", nil, "asset")
	checking, _ := acctSvc.CreateAccount("Checking", &assets.ID, "")
	expenses, _ := acctSvc.CreateAccount("Expenses", nil, "expense")

	// Create a transaction first
	body, _ := json.Marshal(map[string]any{
		"date":        "2026-01-15",
		"description": "original",
		"status":      "pending",
		"postings": []map[string]any{
			{"account_id": checking.ID, "currency": "USD", "amount": "-10.00"},
			{"account_id": expenses.ID, "currency": "USD", "amount": "10.00"},
		},
	})
	createResp, err := http.Post(server.URL+"/api/transactions", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	defer createResp.Body.Close()
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", createResp.StatusCode)
	}
	var created map[string]any
	json.NewDecoder(createResp.Body).Decode(&created)
	txID := int(created["id"].(float64))

	// Regression: PATCH must return 404/405 (not registered), PUT must return 200.
	updateBody, _ := json.Marshal(map[string]any{
		"date":        "2026-01-16",
		"description": "updated",
		"status":      "cleared",
		"postings": []map[string]any{
			{"account_id": checking.ID, "currency": "USD", "amount": "-20.00"},
			{"account_id": expenses.ID, "currency": "USD", "amount": "20.00"},
		},
	})

	// PATCH should NOT work (not registered)
	patchReq, _ := http.NewRequest("PATCH", fmt.Sprintf("%s/api/transactions/%d", server.URL, txID), bytes.NewBuffer(updateBody))
	patchReq.Header.Set("Content-Type", "application/json")
	patchResp, err := http.DefaultClient.Do(patchReq)
	if err != nil {
		t.Fatalf("patch request failed: %v", err)
	}
	defer patchResp.Body.Close()
	if patchResp.StatusCode == http.StatusOK {
		t.Error("PATCH /api/transactions/{id} should not return 200 — only PUT is registered")
	}

	// PUT must work and return 200
	putReq, _ := http.NewRequest("PUT", fmt.Sprintf("%s/api/transactions/%d", server.URL, txID), bytes.NewBuffer(updateBody))
	putReq.Header.Set("Content-Type", "application/json")
	putResp, err := http.DefaultClient.Do(putReq)
	if err != nil {
		t.Fatalf("put request failed: %v", err)
	}
	defer putResp.Body.Close()
	if putResp.StatusCode != http.StatusOK {
		t.Errorf("PUT /api/transactions/{id} expected 200, got %d", putResp.StatusCode)
	}
}
