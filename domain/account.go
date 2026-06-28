package domain

import (
	"fmt"
	"time"
)

type Type string

const (
	Asset     Type = "Asset"
	Liability Type = "Liability"
	Revenue   Type = "Revenue"
	Expense   Type = "Expense"
	Equity    Type = "Equity"
)

type Account struct {
	id          int64
	parentID    int64
	name        string
	type_       Type
	currency    string
	description string
	active      bool
	dateCreated time.Time
	dateUpdated time.Time
}

type AccountDTO struct {
	Id          int64  `json:"id"`
	ParentID    int64  `json:"parentID"`
	Name        string `json:"name"`
	Type_       Type   `json:"type"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

func (a *AccountDTO) ToDomain() (*Account, error) {
	created, err := time.Parse(time.RFC3339, a.DateCreated)
	if err != nil {
		return nil, err
	}
	updated, err := time.Parse(time.RFC3339, a.DateUpdated)
	if err != nil {
		return nil, err
	}

	return NewAccount(
		a.Id,
		a.ParentID,
		a.Name,
		a.Type_,
		a.Currency,
		a.Description,
		a.Active,
		created.Local(),
		updated.Local(),
	), nil
}

func (a *Account) ToDTO() *AccountDTO {
	return &AccountDTO{
		Id:          a.id,
		ParentID:    a.parentID,
		Name:        a.name,
		Type_:       a.type_,
		Currency:    a.currency,
		Description: a.description,
		Active:      a.active,
		DateCreated: a.dateCreated.Format("2006-01-02 15:04:05"),
		DateUpdated: a.dateUpdated.Format("2006-01-02 15:04:05"),
	}
}

func NewAccount(
	id int64,
	parentID int64,
	name string,
	type_ Type,
	currency string,
	description string,
	active bool,
	dateCreated time.Time,
	dateUpdated time.Time,
) *Account {
	a := &Account{
		id,
		parentID,
		name,
		type_,
		currency,
		description,
		active,
		dateCreated,
		dateUpdated,
	}
	return a
}

func (a *Account) ID() int64              { return a.id }
func (a *Account) ParentID() int64        { return a.parentID }
func (a *Account) Name() string           { return a.name }
func (a *Account) Type() Type             { return a.type_ }
func (a *Account) Currency() string       { return a.currency }
func (a *Account) Description() string    { return a.description }
func (a *Account) Active() bool           { return a.active }
func (a *Account) DateCreated() time.Time { return a.dateCreated }
func (a *Account) DateUpdated() time.Time { return a.dateUpdated }

// Getters, Setters that call the DB functions to edit

// move account tree to a new parent
func (a *Account) SetParentID(newParendID int64) error {
	a.parentID = newParendID
	a.update()
	return nil
}

// change Account name
func (a *Account) SetName(newName string) error {
	if err := validateName(newName); err != nil {
		return err
	}
	a.name = newName
	a.update()
	return nil
}

// change Account type
func (a *Account) SetType(newType Type) error {
	a.type_ = newType
	a.update()
	return nil
}

// change description
func (a *Account) SetDescription(newDesc string) error {
	a.description = newDesc
	a.update()
	return nil
}

// set active or inactive
func (a *Account) SetActive(active bool) error {
	a.active = active
	a.update()
	return nil
}

func (a *Account) update() {
	a.dateUpdated = time.Now()
}

func validateName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("Validation error: Account name must exist")
	}
	return nil
}

func (a *Account) Validate() error {
	if err := validateName(a.name); err != nil {
		return err
	}
	return nil
}
