package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/crispy-2/backend/internal/model"
	"github.com/crispy-2/backend/internal/service"
	"github.com/crispy-2/backend/internal/testutil"
)

// wirePagedResult matches the {data: [...], meta: {...}} wire format.
type wirePagedResult[T any] struct {
	Data []T `json:"data"`
	Meta struct {
		Total      int `json:"total"`
		TotalPages int `json:"total_pages"`
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
	} `json:"meta"`
}

func setupAccountRestTest(t *testing.T) (*service.AccountService, *httptest.Server) {
	t.Helper()
	db := testutil.NewTestDB(t)
	svc := service.NewAccountService(db)
	mux := SetupRouter(db, svc, nil, nil, nil, nil, nil, nil, nil, nil)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	return svc, server
}

func TestCreateAccountSuccess(t *testing.T) {
	_, server := setupAccountRestTest(t)
	defer server.Close()

	reqBody := map[string]any{
		"name":        "Assets",
		"parent_id":   nil,
		"account_type": "asset",
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", server.URL+"/api/accounts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var account model.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if account.Name != "Assets" {
		t.Errorf("expected name 'Assets', got %q", account.Name)
	}
	if account.AccountType != "asset" {
		t.Errorf("expected account_type 'asset', got %q", account.AccountType)
	}
	if account.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestCreateAccountChild(t *testing.T) {
	svc, server := setupAccountRestTest(t)
	defer server.Close()

	parent, err := svc.CreateAccount("Assets", nil, "asset")
	if err != nil {
		t.Fatalf("failed to create parent: %v", err)
	}

	reqBody := map[string]any{
		"name":        "Checking",
		"parent_id":   parent.ID,
		"account_type": "",
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", server.URL+"/api/accounts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var account model.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if account.FullName != "Assets:Checking" {
		t.Errorf("expected full_name 'Assets:Checking', got %q", account.FullName)
	}
	if account.ParentID == nil || *account.ParentID != parent.ID {
		t.Errorf("expected parent_id %d, got %v", parent.ID, account.ParentID)
	}
}

// Story 2: create subaccount via ":" path in name
func TestCreateAccountByPath_REST(t *testing.T) {
	_, server := setupAccountRestTest(t)
	defer server.Close()

	reqBody := map[string]any{
		"name":         "Expenses:Food:Restaurants",
		"account_type": "",
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", server.URL+"/api/accounts", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected status 201, got %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var account model.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if account.FullName != "Expenses:Food:Restaurants" {
		t.Errorf("expected full_name 'Expenses:Food:Restaurants', got %q", account.FullName)
	}
	if account.Name != "Restaurants" {
		t.Errorf("expected name 'Restaurants', got %q", account.Name)
	}
}

func TestCreateAccountValidationErrors(t *testing.T) {
	tests := []struct {
		name          string
		reqBody       map[string]any
		expectedCode  int
		expectedError string
	}{
		{
			name: "empty name",
			reqBody: map[string]any{
				"name":        "",
				"parent_id":   nil,
				"account_type": "asset",
			},
			expectedCode:  400,
			expectedError: "account name must not be empty",
		},
		{
			name: "invalid account type",
			reqBody: map[string]any{
				"name":        "Test",
				"parent_id":   nil,
				"account_type": "invalid_type",
			},
			expectedCode:  400,
			expectedError: "invalid account_type: invalid_type",
		},
		{
			name: "root without type",
			reqBody: map[string]any{
				"name":        "Test",
				"parent_id":   nil,
				"account_type": "",
			},
			expectedCode:  400,
			expectedError: "root accounts must have an account_type",
		},
		{
			name: "child with parent but wrong parent type",
			reqBody: map[string]any{
				"name":        "Child",
				"parent_id":   int64(999),
				"account_type": "",
			},
			expectedCode:  404,
			expectedError: "not found",
		},
		{
			name: "name contains separator — creates hierarchy",
			reqBody: map[string]any{
				"name":        "Bad:Name",
				"parent_id":   nil,
				"account_type": "asset",
			},
			// With ":" in name, the handler uses CreateAccountByPath which creates the hierarchy.
			// This is now valid — skip from validation error tests.
			expectedCode:  0, // sentinel: skip this case
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedCode == 0 {
				t.Skip("case intentionally removed — behaviour changed")
			}
			_, server := setupAccountRestTest(t)
			defer server.Close()

			body, _ := json.Marshal(tc.reqBody)
			req, err := http.NewRequest("POST", server.URL+"/api/accounts", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("failed to execute request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("expected status %d, got %d", tc.expectedCode, resp.StatusCode)
			}

			var appErr model.AppError
			if err := json.NewDecoder(resp.Body).Decode(&appErr); err != nil {
				t.Fatalf("failed to decode error response: %v", err)
			}

			if appErr.Code != tc.expectedCode {
				t.Errorf("expected error code %d, got %d", tc.expectedCode, appErr.Code)
			}
		})
	}
}

func TestGetAccountSuccess(t *testing.T) {
	svc, server := setupAccountRestTest(t)
	defer server.Close()

	created, err := svc.CreateAccount("Assets", nil, "asset")
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	req, err := http.NewRequest("GET", server.URL+"/api/accounts/"+fmt.Sprintf("%d", created.ID), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var account model.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if account.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, account.ID)
	}
	if account.Name != "Assets" {
		t.Errorf("expected name 'Assets', got %q", account.Name)
	}
}

func TestGetAccountNotFound(t *testing.T) {
	_, server := setupAccountRestTest(t)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL+"/api/accounts/9999", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	var appErr model.AppError
	if err := json.NewDecoder(resp.Body).Decode(&appErr); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if appErr.Code != 404 {
		t.Errorf("expected error code 404, got %d", appErr.Code)
	}
}

func TestListAccounts(t *testing.T) {
	svc, server := setupAccountRestTest(t)
	defer server.Close()

	svc.CreateAccount("Assets", nil, "asset")
	svc.CreateAccount("Expenses", nil, "expense")
	svc.CreateAccount("Liabilities", nil, "liability")

	req, err := http.NewRequest("GET", server.URL+"/api/accounts?page=1&page_size=10", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var result wirePagedResult[model.Account]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Meta.Total != 3 {
		t.Errorf("expected 3 accounts, got %d", result.Meta.Total)
	}
	if len(result.Data) != 3 {
		t.Errorf("expected 3 items, got %d", len(result.Data))
	}
}

func TestListAccountsFilterByType(t *testing.T) {
	svc, server := setupAccountRestTest(t)
	defer server.Close()

	svc.CreateAccount("Assets", nil, "asset")
	svc.CreateAccount("Expenses", nil, "expense")
	svc.CreateAccount("Liabilities", nil, "liability")

	req, err := http.NewRequest("GET", server.URL+"/api/accounts?account_type=asset", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	var result wirePagedResult[model.Account]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Meta.Total != 1 {
		t.Errorf("expected 1 account, got %d", result.Meta.Total)
	}
	if len(result.Data) > 0 && result.Data[0].AccountType != "asset" {
		t.Errorf("expected asset type, got %q", result.Data[0].AccountType)
	}
}

func TestUpdateAccountSuccess(t *testing.T) {
	svc, server := setupAccountRestTest(t)
	defer server.Close()

	created, err := svc.CreateAccount("Assets", nil, "asset")
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	reqBody := map[string]any{
		"name": "UpdatedAssets",
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("PUT", server.URL+"/api/accounts/"+fmt.Sprintf("%d", created.ID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var account model.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if account.Name != "UpdatedAssets" {
		t.Errorf("expected name 'UpdatedAssets', got %q", account.Name)
	}
}

func TestDeleteAccountSuccess(t *testing.T) {
	svc, server := setupAccountRestTest(t)
	defer server.Close()

	created, err := svc.CreateAccount("Temp", nil, "asset")
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}
	accountID := created.ID

	req, err := http.NewRequest("DELETE", server.URL+"/api/accounts/"+fmt.Sprintf("%d", accountID), nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	_, err = svc.GetAccount(accountID)
	if err == nil {
		t.Error("expected account to be deleted")
	}
}

func TestDeleteAccountNotFound(t *testing.T) {
	_, server := setupAccountRestTest(t)
	defer server.Close()

	req, err := http.NewRequest("DELETE", server.URL+"/api/accounts/9999", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestSearchAccountsSuccess(t *testing.T) {
	svc, server := setupAccountRestTest(t)
	defer server.Close()

	svc.CreateAccount("Assets", nil, "asset")
	svc.CreateAccount("Expenses", nil, "expense")
	svc.CreateAccount("Liabilities", nil, "liability")

	req, err := http.NewRequest("GET", server.URL+"/api/accounts/search?q=asset", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to execute request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var accounts []model.Account
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("expected 1 account, got %d", len(accounts))
	}
	if len(accounts) > 0 && accounts[0].Name != "Assets" {
		t.Errorf("expected 'Assets', got %q", accounts[0].Name)
	}
}
