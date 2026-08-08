package domain

import (
	"time"
)

type Budget struct {
	id          int64
	name        string
	description string
	target      string
	tagId       int64
	recurrenceId int64
}

type BudgetDTO struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Target       string `json:"target"`
	TagId        int64  `json:"tagId"`
	TagName      string `json:"tagName"`
	Frequency    string `json:"frequency"`
	Interval     int    `json:"interval"`
	DayOfWeek    int    `json:"dayOfWeek"`
	DayOfMonth   int    `json:"dayOfMonth"`
	StartDate    string `json:"startDate"`
	CurrentSpend string `json:"currentSpend"`
	NextReset    string `json:"nextReset"`
}

func NewBudget(
	id int64,
	name string,
	description string,
	target string,
	tagId int64,
	recurrenceId int64,
) *Budget {
	return &Budget{
		id:           id,
		name:         name,
		description:  description,
		target:       target,
		tagId:        tagId,
		recurrenceId: recurrenceId,
	}
}

func (b *Budget) ToDTO() *BudgetDTO {
	return &BudgetDTO{
		Id:          b.id,
		Name:        b.name,
		Description: b.description,
		Target:      b.target,
		TagId:       b.tagId,
	}
}

func (b *BudgetDTO) ToDomain() (*Budget, error) {
	return NewBudget(
		b.Id,
		b.Name,
		b.Description,
		b.Target,
		b.TagId,
		0, // recurrenceId resolved by service layer
	), nil
}

func (b *Budget) ID() int64           { return b.id }
func (b *Budget) Name() string        { return b.name }
func (b *Budget) Description() string { return b.description }
func (b *Budget) Target() string      { return b.target }
func (b *Budget) TagId() int64        { return b.tagId }
func (b *Budget) RecurrenceId() int64 { return b.recurrenceId }

// RecurrenceInput carries the recurrence fields from the frontend DTO
// so the service layer can upsert the Recurrence row independently.
type RecurrenceInput struct {
	Frequency  string
	Interval   int
	DayOfWeek  time.Weekday
	DayOfMonth int
	StartDate  time.Time
}

func (b *BudgetDTO) ToRecurrenceInput() (*RecurrenceInput, error) {
	startDate, err := time.Parse("2006-01-02", b.StartDate)
	if err != nil {
		return nil, err
	}
	return &RecurrenceInput{
		Frequency:  b.Frequency,
		Interval:   b.Interval,
		DayOfWeek:  time.Weekday(b.DayOfWeek),
		DayOfMonth: b.DayOfMonth,
		StartDate:  startDate,
	}, nil
}
