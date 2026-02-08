package domain

import "time"

type RecurringTransaction struct {
	id          int64
	description string
	schedule    Recurrence
	lastRun     time.Time
	active      bool
	dateCreated time.Time
	dateUpdated time.Time
}

func NewRecurringTransaction(
	id int64,
	description string,
	schedule Recurrence,
	lastRun time.Time,
	active bool,
	dateCreated time.Time,
	dateUpdated time.Time,
) *RecurringTransaction {
	return &RecurringTransaction{
		id,
		description,
		schedule,
		lastRun,
		active,
		dateCreated,
		dateUpdated,
	}
}

func (r *RecurringTransaction) ID() int64              { return r.id }
func (r *RecurringTransaction) Description() string    { return r.description }
func (r *RecurringTransaction) Schedule() Recurrence   { return r.schedule }
func (r *RecurringTransaction) LastRun() time.Time     { return r.lastRun }
func (r *RecurringTransaction) Active() bool           { return r.active }
func (r *RecurringTransaction) DateCreated() time.Time { return r.dateCreated }
func (r *RecurringTransaction) DateUpdated() time.Time { return r.dateUpdated }
