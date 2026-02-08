package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type Status string

const (
	Pending Status = "Pending"
	Posted  Status = "Posted"
	Cleared Status = "Cleared"
)

type Transaction struct {
	id          int64
	description string
	date        time.Time
	status      Status
	referenceID *int64
	dateCreated time.Time
	dateUpdated time.Time
	postings    []*Posting
	tags        []string
}

type TransactionDTO struct {
	Id          int64         `json:"id"`
	Description string        `json:"description"`
	Date        string        `json:"date"`
	Status      Status        `json:"status"`
	ReferenceID *int64        `json:"referenceID"`
	DateCreated string        `json:"dateCreated"`
	DateUpdated string        `json:"dateUpdated"`
	Postings    []*PostingDTO `json:"postings"`
	Tags        []string      `json:"tags"`
}

func (t *TransactionDTO) ToDomain() (*Transaction, error) {
	converted, err := time.Parse(time.RFC3339, t.Date)
	if err != nil {
		return nil, err
	}
	localTime := converted.Local()
	y, m, d := localTime.Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, localTime.Location())
	created, err := time.Parse(time.RFC3339, t.DateCreated)
	if err != nil {
		return nil, err
	}
	updated, err := time.Parse(time.RFC3339, t.DateUpdated)
	if err != nil {
		return nil, err
	}
	var postings []*Posting
	for _, p := range t.Postings {
		posting, err := p.ToDomain()
		if err != nil {
			return nil, err
		}
		postings = append(postings, posting)
	}
	return NewTransaction(
		t.Id,
		t.Description,
		date,
		t.Status,
		t.ReferenceID,
		created.Local(),
		updated.Local(),
		postings,
		t.Tags,
	)
}

func (t *Transaction) ToDTO() *TransactionDTO {
	var postings []*PostingDTO
	for _, p := range t.postings {
		postings = append(postings, p.ToDTO())
	}
	return &TransactionDTO{
		Id:          t.id,
		Description: t.description,
		Date:        t.date.Format(time.RFC3339),
		Status:      t.status,
		ReferenceID: t.referenceID,
		DateCreated: t.dateCreated.Format("2006-01-02 15:04:05"),
		DateUpdated: t.dateUpdated.Format("2006-01-02 15:04:05"),
		Postings:    postings,
		Tags:        t.tags,
	}
}

func NewTransaction(
	id int64,
	description string,
	date time.Time,
	status Status,
	referenceID *int64,
	dateCreated time.Time,
	dateUpdated time.Time,
	postings []*Posting,
	tags []string,
) (*Transaction, error) {
	tx := &Transaction{
		id,
		description,
		date,
		status,
		referenceID,
		dateCreated,
		dateUpdated,
		postings,
		tags,
	}
	if err := tx.Validate(); err != nil {
		return nil, err
	}
	return tx, nil
}

func (t *TransactionDTO) String() string {
	var b strings.Builder
	tokens := []struct {
		f string
		a any
	}{
		{"ID: %d", t.Id},
		{"Desc: %s", t.Description},
		{"Date: %s", t.Date},
		{"Status: %s", t.Status},
		{"RefID: %d", t.ReferenceID},
		{"Created: %s", t.DateCreated},
		{"Updated: %s", t.DateUpdated},
		{"Postings: %v", t.Postings},
		{"Tags: %s", t.Tags},
	}
	for i, s := range tokens {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf(s.f, s.a))
	}
	return b.String()
}

func (t *Transaction) String() string {
	var b strings.Builder
	tokens := []struct {
		f string
		a any
	}{
		{"ID: %d", t.id},
		{"Desc: %s", t.description},
		{"Date: %s", t.date},
		{"Status: %s", t.status},
		{"RefID: %d", t.referenceID},
		{"Created: %s", t.dateCreated},
		{"Updated: %s", t.dateUpdated},
		{"Postings: %v", t.postings},
		{"Tags: %s", t.tags},
	}
	for i, s := range tokens {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(fmt.Sprintf(s.f, s.a))
	}
	return b.String()
}

func (t *Transaction) ID() int64              { return t.id }
func (t *Transaction) Description() string    { return t.description }
func (t *Transaction) Date() time.Time        { return t.date }
func (t *Transaction) Status() Status         { return t.status }
func (t *Transaction) ReferenceID() *int64    { return t.referenceID }
func (t *Transaction) DateCreated() time.Time { return t.dateCreated }
func (t *Transaction) DateUpdated() time.Time { return t.dateUpdated }
func (t *Transaction) Postings() []*Posting   { return t.postings }
func (t *Transaction) Tags() []string         { return t.tags }

func (t *Transaction) SetDescription(newDesc string) error {
	t.description = newDesc
	t.update()
	return nil
}

func (t *Transaction) SetDate(newDate time.Time) error {
	t.date = newDate
	t.update()
	return nil
}

func (t *Transaction) SetStatus(newStatus Status) error {
	t.status = newStatus
	t.update()
	return nil
}

func (t *Transaction) SetReferenceID(newID *int64) error {
	t.referenceID = newID
	t.update()
	return nil
}

func (t *Transaction) SetPostings(newPostings []*Posting) error {
	if err := validatePostings(newPostings); err != nil {
		return err
	}
	t.postings = newPostings
	t.update()
	return nil
}

func (t *Transaction) SetTags(newTags []string) {
	t.tags = newTags
	t.update()
}

func (t *Transaction) update() {
	t.dateUpdated = time.Now()
}

func validatePostings(postings []*Posting) error {
	if len(postings) == 0 {
		return fmt.Errorf("Validation error: Transaction must contain postings")
	}
	var total decimal.Decimal
	for _, p := range postings {
		total.Add(p.amount)
	}
	if !total.Equal(decimal.Zero) {
		return fmt.Errorf("Validation error: Unbalanced transaction")
	}
	return nil
}

func (t *Transaction) Validate() error {
	if err := validatePostings(t.postings); err != nil {
		return err
	}
	return nil
}
