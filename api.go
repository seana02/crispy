package main

import (
	"context"
	"crispy/domain"
	"fmt"
)

func (a *App) CreateTransactionSample(t *domain.TransactionDTO) int {
	domainObj, err := t.ToDomain()
	if err != nil {
		fmt.Printf("%s\n", err)
		return -1
	}
	fmt.Printf("%v\f", domainObj)
	return 19
}

func (a *App) CreateAccountSample(t *domain.AccountDTO) int {
	accountObj, err := t.ToDomain()
	if err != nil {
		fmt.Printf("%s\n", err)
		return -1
	}
	fmt.Printf("%v\f", accountObj)
	return 19
}

func (a *App) CreateTransaction(t *domain.TransactionDTO) (*domain.TransactionDTO, error) {
	return a.SubmitTransaction(a.deps.Service.CreateTransaction, t)
}

func (a *App) UpdateTransaction(t *domain.TransactionDTO) (*domain.TransactionDTO, error) {
	return a.SubmitTransaction(a.deps.Service.UpdateTransaction, t)
}

func (a *App) SubmitTransaction(fn func(ctx context.Context, t *domain.Transaction) (*domain.Transaction, error), t *domain.TransactionDTO) (*domain.TransactionDTO, error) {
	domainObj, err := t.ToDomain()
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	returnedTransaction, err := fn(a.ctx, domainObj)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	return returnedTransaction.ToDTO(), nil
}

func (a *App) CreateAccount(t *domain.AccountDTO) (*domain.AccountDTO, error) {
	acctObj, err := t.ToDomain()
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	returnedAccount, err := a.deps.Service.CreateAccount(a.ctx, acctObj)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	return returnedAccount.ToDTO(), nil
}

func (a *App) GetAccountID(acct string) int64 {
	id, err := a.deps.Service.GetID(a.ctx, acct)
	if err != nil {
		fmt.Printf("%s\n", err)
		return -1
	}
	return id
}

func (a *App) GetTransactionList() ([]*domain.TransactionDTO, error) {
	t_list, err := a.deps.Service.GetAllTransactions(a.ctx)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	output := []*domain.TransactionDTO{}
	for _, t := range t_list {
		newT, err := a.prepareTransaction(t)
		if err != nil {
			fmt.Printf("%s\n", err)
			return nil, err
		}
		output = append(output, newT)
	}
	return output, nil
}

func (a *App) GetTransactionByID(id int64) (*domain.TransactionDTO, error) {
	tx, err := a.deps.Service.GetTransactionByID(a.ctx, id)
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	return a.prepareTransaction(tx)
}

func (a *App) prepareTransaction(tx *domain.Transaction) (*domain.TransactionDTO, error) {
	t := tx.ToDTO()
	for _, p := range t.Postings {
		name, err := a.deps.Service.GetAccountName(a.ctx, p.AccountID)
		if err != nil {
			fmt.Printf("%s\n", err)
			return nil, err
		}
		p.AccountName = name
	}
	return t, nil
}
