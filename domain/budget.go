package domain

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type Frequency string

const (
	FreqDaily   Frequency = "Daily"
	FreqWeekly  Frequency = "Weekly"
	FreqMonthly Frequency = "Monthly"
)

type Budget struct {
	id          int64
	name        string
	description string
	target      decimal.Decimal
	tagId       int64
	frequency   Frequency
	interval    int
	dayOfWeek   *time.Weekday
	dayOfMonth  int
}

type BudgetDTO struct {
	Id           int64           `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Target       decimal.Decimal `json:"target"`
	TagId        int64           `json:"tagId"`
	Frequency    Frequency       `json:"frequency"`
	Interval     int             `json:"interval"`
	DayOfWeek    *time.Weekday   `json:"dayOfWeek"`
	DayOfMonth   int             `json:"dayOfMonth"`
	CurrentSpend string          `json:"currentSpend"`
	NextReset    string          `json:"nextReset"`
}

func NewBudget(
	id int64,
	name string,
	description string,
	target decimal.Decimal,
	tagId int64,
	frequency Frequency,
	interval int,
	dayOfWeek *time.Weekday,
	dayOfMonth int,
) *Budget {
	return &Budget{
		id:          id,
		name:        name,
		description: description,
		target:      target,
		tagId:       tagId,
		frequency:   frequency,
		interval:    interval,
		dayOfWeek:   dayOfWeek,
		dayOfMonth:  dayOfMonth,
	}
}

func (b *Budget) ToDTO() *BudgetDTO {
	return &BudgetDTO{
		Id:          b.id,
		Name:        b.name,
		Description: b.description,
		Target:      b.target,
		TagId:       b.tagId,
		Frequency:   b.frequency,
		Interval:    b.interval,
		DayOfWeek:   b.dayOfWeek,
		DayOfMonth:  b.dayOfMonth,
	}
}

func (b *BudgetDTO) ToDomain() (*Budget, error) {
	return NewBudget(
		b.Id,
		b.Name,
		b.Description,
		b.Target,
		b.TagId,
		b.Frequency,
		b.Interval,
		b.DayOfWeek,
		b.DayOfMonth,
	), nil
}

func (b *Budget) ID() int64                { return b.id }
func (b *Budget) Name() string             { return b.name }
func (b *Budget) Description() string      { return b.description }
func (b *Budget) Target() decimal.Decimal  { return b.target }
func (b *Budget) TagId() int64             { return b.tagId }
func (b *Budget) Frequency() Frequency     { return b.frequency }
func (b *Budget) Interval() int            { return b.interval }
func (b *Budget) DayOfWeek() *time.Weekday { return b.dayOfWeek }
func (b *Budget) DayOfMonth() int          { return b.dayOfMonth }

func (b *Budget) Validate() error {
	if b.frequency == FreqWeekly && b.dayOfWeek == nil {
		return errors.New("Weekly frequency requires a day of week")
	}
	if b.frequency == FreqMonthly && (b.dayOfMonth < 0 || b.dayOfMonth > 31) {
		return errors.New("Monthly frequency requires a day of month")
	}
	return nil
}
