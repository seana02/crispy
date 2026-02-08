package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Posting struct {
	id            int64
	transactionID int64
	accountID     int64
	amount        decimal.Decimal
	currency      string
	dateCreated   time.Time
	dateUpdated   time.Time
}

type PostingDTO struct {
	Id            int64           `json:"id"`
	TransactionID int64           `json:"transactionID"`
	AccountID     int64           `json:"accountID"`
	AccountName   string          `json:"accountName"`
	Amount        decimal.Decimal `json:"amount"`
	Currency      string          `json:"currency"`
	DateCreated   string          `json:"dateCreated"`
	DateUpdated   string          `json:"dateUpdated"`
}

func (p *PostingDTO) String() string {
	var b strings.Builder
	tokens := []struct {
		f string
		a any
	}{
		{"ID: %d", p.Id},
		{"AccountID: %d", p.AccountID},
		{"AccountName: %s", p.AccountName},
		{"Amount: %s", p.Amount},
		{"Currency: %s", p.Currency},
		{"Created: %s", p.DateCreated},
		{"Updated: %s", p.DateUpdated},
	}
	for i, s := range tokens {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf(s.f, s.a))
	}
	return b.String()
}

func (p *Posting) String() string {
	var b strings.Builder
	tokens := []struct {
		f string
		a any
	}{
		{"ID: %d", p.id},
		{"AccountID: %d", p.accountID},
		{"Amount: %s", p.amount},
		{"Currency: %s", p.currency},
		{"Created: %s", p.dateCreated},
		{"Updated: %s", p.dateUpdated},
	}
	for i, s := range tokens {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf(s.f, s.a))
	}
	return b.String()
}

func (p *PostingDTO) ToDomain() (*Posting, error) {
	created, err := time.Parse(time.RFC3339, p.DateCreated)
	if err != nil {
		return nil, err
	}
	updated, err := time.Parse(time.RFC3339, p.DateCreated)
	if err != nil {
		return nil, err
	}
	return NewPosting(
		p.Id,
		p.TransactionID,
		p.AccountID,
		p.Amount,
		p.Currency,
		created.Local(),
		updated.Local(),
	), nil
}

func (p *Posting) ToDTO() *PostingDTO {
	return &PostingDTO{
		Id:            p.ID(),
		TransactionID: p.TransactionID(),
		AccountID:     p.AccountID(),
		Amount:        p.Amount(),
		Currency:      p.Currency(),
		DateCreated:   p.DateCreated().Format(time.RFC3339),
		DateUpdated:   p.DateUpdated().Format(time.RFC3339),
	}
}

func NewPosting(
	id int64,
	transactionID int64,
	accountID int64,
	amount decimal.Decimal,
	currency string,
	dateCreated time.Time,
	dateUpdated time.Time,
) *Posting {
	return &Posting{
		id,
		transactionID,
		accountID,
		amount,
		currency,
		dateCreated,
		dateUpdated,
	}
}

func (p *Posting) ID() int64               { return p.id }
func (p *Posting) TransactionID() int64    { return p.transactionID }
func (p *Posting) AccountID() int64        { return p.accountID }
func (p *Posting) Amount() decimal.Decimal { return p.amount }
func (p *Posting) Currency() string        { return p.currency }
func (p *Posting) DateCreated() time.Time  { return p.dateCreated }
func (p *Posting) DateUpdated() time.Time  { return p.dateUpdated }

func (p *Posting) SetTransactionID(newTransactionID int64) {
	if p.transactionID != -1 {
		return
	}
	p.transactionID = newTransactionID
	p.update()
}

// change associated Account
func (p *Posting) SetAccount(newAccountID int64) {
	p.accountID = newAccountID
	p.update()
}

// change amount
func (p *Posting) SetAmount(newAmount decimal.Decimal) {
	p.amount = newAmount
	p.update()
}

// change currency
func (p *Posting) SetCurrency(newCurrency string) {
	p.currency = newCurrency
	p.update()
}

func (p *Posting) update() {
	p.dateUpdated = time.Now()
}
