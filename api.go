package main

import (
	"context"
	"crispy/domain"
	"fmt"
)

// =============================================================== //
// *                                                             * //
// *                        Transactions                         * //
// *                                                             * //
// =============================================================== //

func (a *App) CreateTransaction(t *domain.TransactionDTO) error {
	if err := a.SubmitTransaction(a.deps.Service.CreateTransaction, t); err != nil {
		a.deps.Logger.Error("Error creating transaction", "Error", err)
		return err
	}
	return nil
}

func (a *App) UpdateTransaction(t *domain.TransactionDTO) error {
	if err := a.SubmitTransaction(a.deps.Service.UpdateTransaction, t); err != nil {
		a.deps.Logger.Error("Error updating transaction", "Error", err)
		return err
	}
	return nil
}

func (a *App) SubmitTransaction(fn func(ctx context.Context, t *domain.Transaction) error, t *domain.TransactionDTO) error {
	domainObj, err := t.ToDomain()
	if err != nil {
		return err
	}
	err = fn(a.ctx, domainObj)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) GetTransactionList() ([]*domain.TransactionDTO, error) {
	t_list, err := a.deps.Service.GetAllTransactions(a.ctx, -1, -1)
	if err != nil {
		a.deps.Logger.Error("Error getting all transactions", "Error", err)
		return nil, err
	}
	output := []*domain.TransactionDTO{}
	for _, t := range t_list {
		newT, err := a.prepareTransaction(t)
		if err != nil {
			return nil, err
		}
		output = append(output, newT)
	}
	return output, nil
}

func (a *App) GetTransactionByID(id int64) (*domain.TransactionDTO, error) {
	tx, err := a.deps.Service.GetTransactionByID(a.ctx, id)
	if err != nil {
		a.deps.Logger.Error("Error getting transaction", "id", id, "Error", err)
		return nil, err
	}
	return a.prepareTransaction(tx)
}

func (a *App) prepareTransaction(tx *domain.Transaction) (*domain.TransactionDTO, error) {
	t := tx.ToDTO()
	for _, p := range t.Postings {
		name, err := a.deps.Service.GetAccountName(a.ctx, p.AccountID)
		if err != nil {
			a.deps.Logger.Error("Error preparing transaction", "Error", err)
			return nil, err
		}
		p.AccountName = name
	}
	return t, nil
}

func (a *App) DeleteTransaction(id int64) error {
	if err := a.deps.Service.DeleteTransaction(a.ctx, id); err != nil {
		a.deps.Logger.Error("Error deleting transaction", "id", id, "Error", err)
	}
	return nil
}

// =============================================================== //4/10
// *                                                             * //
// *                          Accounts                           * //
// *                                                             * //
// =============================================================== //

func (a *App) CreateAccount(acc *domain.AccountDTO) error {
	if err := a.SubmitAccount(a.deps.Service.CreateAccount, acc); err != nil {
		a.deps.Logger.Error("Error creating account", "Error", err)
		return err
	}
	return nil
}

func (a *App) UpdateAccount(acc *domain.AccountDTO) error {
	if err := a.SubmitAccount(a.deps.Service.UpdateAccount, acc); err != nil {
		a.deps.Logger.Error("Error updating account", "Error", err)
		return err
	}
	return nil
}

func (a *App) SubmitAccount(fn func(context.Context, *domain.Account) error, t *domain.AccountDTO) error {
	acctObj, err := t.ToDomain()
	if err != nil {
		return err
	}
	err = fn(a.ctx, acctObj)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) GetAccountList() ([]*domain.AccountDTO, error) {
	accts, err := a.deps.Service.GetAllAccounts(a.ctx)
	if err != nil {
		a.deps.Logger.Error("Error getting account list", "Error", err)
		return nil, err
	}
	var output []*domain.AccountDTO
	for _, acc := range accts {
		output = append(output, acc.ToDTO())
	}
	return output, nil
}

func (a *App) GetAccountIDByName(acct string) int64 {
	id, err := a.deps.Service.GetAccountID(a.ctx, acct)
	if err != nil {
		a.deps.Logger.Error("Error getting account id", "account name", acct, "Error", err)
		return -1
	}
	return id
}

func (a *App) GetAccountByID(id int64) (*domain.AccountDTO, error) {
	acct, err := a.deps.Service.GetAccountByID(a.ctx, id)
	if err != nil {
		a.deps.Logger.Error("Error getting account", "id", id, "Error", err)
		return nil, err
	}
	return acct.ToDTO(), nil
}

func (a *App) DeleteAccount(id int64) error {
	if err := a.deps.Service.DeleteAccount(a.ctx, id); err != nil {
		a.deps.Logger.Error("Error deleting account", "id", id, "Error", err)
	}
	return nil
}

// =============================================================== //
// *                                                             * //
// *                 Sample Functions for Testing                * //
// *                                                             * //
// =============================================================== //

func (a *App) CreateTransactionSample(t *domain.TransactionDTO) int {
	domainObj, err := t.ToDomain()
	if err != nil {
		a.deps.Logger.Error("Error converting from DTO", "Error", err)
		return -1
	}
	a.deps.Logger.Debug("CreateTransactionSample", "Transaction Object", domainObj)
	fmt.Printf("%v\f", domainObj)
	return 19
}

func (a *App) CreateAccountSample(t *domain.AccountDTO) int {
	accountObj, err := t.ToDomain()
	if err != nil {
		a.deps.Logger.Error("Error converting from DTO", "Error", err)
		return -1
	}
	a.deps.Logger.Debug("CreateAccountSample", "Transaction Object", accountObj)
	return 19
}
