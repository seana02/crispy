package service

import (
    "database/sql"
    "fmt"
    "log/slog"
    "time"

    "github.com/crispy-2/backend/internal/dal/generated"
    "github.com/crispy-2/backend/internal/model"
    "github.com/shopspring/decimal"
)

// RecurringService handles recurring transaction templates.
type RecurringService struct {
    db     *sql.DB
    queries *generated.Queries
    txSvc  *TransactionService
}

// NewRecurringService creates a new RecurringService.
func NewRecurringService(db *sql.DB, txSvc *TransactionService) *RecurringService {
    return &RecurringService{
        db:      db,
        queries: generated.New(db),
        txSvc:   txSvc,
    }
}

// CreateTemplate creates a new recurring transaction template.
func (s *RecurringService) CreateTemplate(
    description, status, scheduleType string, intervalN int, dayOfMonth *int,
    startDate string, endDate *string,
    postings []PostingInput, tags []TagInput,
) (*model.RecurringTemplate, error) {
    if err := validateSchedule(scheduleType, intervalN, dayOfMonth); err != nil {
        return nil, err
    }
    if !model.ValidStatus(status) {
        return nil, model.NewValidationError("invalid status: " + status)
    }
    if startDate == "" {
        return nil, model.NewValidationError("start_date is required")
    }
    if len(postings) < 2 {
        return nil, model.NewValidationError("at least two postings are required")
    }

    tx, err := s.db.Begin()
    if err != nil {
        return nil, model.NewInternalError("failed to begin transaction: " + err.Error())
    }
    defer tx.Rollback()

    q := generated.New(tx)

    var domSQL sql.NullInt64
    if dayOfMonth != nil {
        domSQL = sql.NullInt64{Int64: int64(*dayOfMonth), Valid: true}
    }
    var endSQL sql.NullString
    if endDate != nil {
        endSQL = sql.NullString{String: *endDate, Valid: true}
    }

    row, err := q.CreateRecurringTemplate(ctx(), generated.CreateRecurringTemplateParams{
        Description:  description,
        Status:       status,
        ScheduleType: scheduleType,
        IntervalN:    int64(intervalN),
        DayOfMonth:   domSQL,
        StartDate:    startDate,
        EndDate:      endSQL,
        IsActive:     1,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to create recurring template: " + err.Error())
    }

    // Create template postings
    var modelPostings []model.TemplatePosting
    for _, p := range postings {
        amt, err := decimal.NewFromString(p.Amount)
        if err != nil {
            return nil, model.NewValidationError("invalid posting amount: " + p.Amount)
        }
        pRow, err := q.CreateTemplatePosting(ctx(), generated.CreateTemplatePostingParams{
            TemplateID: row.ID,
            AccountID:  p.AccountID,
            Currency:   p.Currency,
            Amount:     amt.String(),
        })
        if err != nil {
            return nil, model.NewInternalError("failed to create template posting: " + err.Error())
        }
        modelPostings = append(modelPostings, model.TemplatePosting{
            ID:         pRow.ID,
            TemplateID: row.ID,
            AccountID:  p.AccountID,
            Currency:   p.Currency,
            Amount:     model.StringDecimal{Decimal: amt},
        })
    }

    // Create/link tags
    var modelTags []model.Tag
    for _, t := range tags {
        tagRow, err := q.CreateTag(ctx(), generated.CreateTagParams{Key: t.Key, Value: t.Value})
        if err != nil {
            return nil, model.NewInternalError("failed to create tag: " + err.Error())
        }
        if err := q.AddTemplateTag(ctx(), generated.AddTemplateTagParams{
            TemplateID: row.ID,
            TagID:      tagRow.ID,
        }); err != nil {
            return nil, model.NewInternalError("failed to link tag: " + err.Error())
        }
        modelTags = append(modelTags, model.Tag{ID: tagRow.ID, Key: t.Key, Value: t.Value})
    }

    if err := tx.Commit(); err != nil {
        return nil, model.NewInternalError("failed to commit: " + err.Error())
    }

    slog.Info("recurring template created", "id", row.ID, "description", description)
    return &model.RecurringTemplate{
        ID:           row.ID,
        Description:  description,
        Status:       status,
        ScheduleType: scheduleType,
        IntervalN:    intervalN,
        DayOfMonth:   dayOfMonth,
        StartDate:    startDate,
        EndDate:      endDate,
        IsActive:     true,
        Postings:     modelPostings,
        Tags:         modelTags,
    }, nil
}

// GetTemplate retrieves a recurring template by ID.
func (s *RecurringService) GetTemplate(id int64) (*model.RecurringTemplate, error) {
    row, err := s.queries.GetRecurringTemplate(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, model.NewNotFoundError("recurring template", id)
        }
        return nil, model.NewInternalError("failed to get template: " + err.Error())
    }

    postingRows, err := s.queries.GetTemplatePostingsWithAccounts(ctx(), id)
    if err != nil {
        return nil, model.NewInternalError("failed to get template postings: " + err.Error())
    }

    tagRows, err := s.queries.GetTemplateTags(ctx(), id)
    if err != nil {
        return nil, model.NewInternalError("failed to get template tags: " + err.Error())
    }

    return toModelTemplateWithAccounts(row, postingRows, tagRows), nil
}

// ListTemplates returns all recurring templates.
func (s *RecurringService) ListTemplates() ([]model.RecurringTemplate, error) {
    rows, err := s.queries.ListRecurringTemplates(ctx())
    if err != nil {
        return nil, model.NewInternalError("failed to list templates: " + err.Error())
    }

    templates := make([]model.RecurringTemplate, len(rows))
    for i, r := range rows {
        tmpl, err := s.GetTemplate(r.ID)
        if err != nil {
            return nil, err
        }
        templates[i] = *tmpl
    }
    return templates, nil
}

// UpdateTemplate updates a recurring template.
func (s *RecurringService) UpdateTemplate(
    id int64, description, status, scheduleType string, intervalN int, dayOfMonth *int,
    startDate string, endDate *string, isActive bool,
    postings []PostingInput, tags []TagInput,
) (*model.RecurringTemplate, error) {
    if err := validateSchedule(scheduleType, intervalN, dayOfMonth); err != nil {
        return nil, err
    }

    _, err := s.queries.GetRecurringTemplate(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, model.NewNotFoundError("recurring template", id)
        }
        return nil, model.NewInternalError("failed to get template: " + err.Error())
    }

    tx, err := s.db.Begin()
    if err != nil {
        return nil, model.NewInternalError("failed to begin transaction: " + err.Error())
    }
    defer tx.Rollback()

    q := generated.New(tx)

    var domSQL sql.NullInt64
    if dayOfMonth != nil {
        domSQL = sql.NullInt64{Int64: int64(*dayOfMonth), Valid: true}
    }
    var endSQL sql.NullString
    if endDate != nil {
        endSQL = sql.NullString{String: *endDate, Valid: true}
    }

    var isActiveInt int64
    if isActive {
        isActiveInt = 1
    }

    _, err = q.UpdateRecurringTemplate(ctx(), generated.UpdateRecurringTemplateParams{
        Description:  description,
        Status:       status,
        ScheduleType: scheduleType,
        IntervalN:    int64(intervalN),
        DayOfMonth:   domSQL,
        StartDate:    startDate,
        EndDate:      endSQL,
        IsActive:     isActiveInt,
        UpdatedAt:    time.Now().Format(time.RFC3339),
        ID:           id,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to update template: " + err.Error())
    }

    // Replace postings
    if err := q.DeleteTemplatePostings(ctx(), id); err != nil {
        return nil, model.NewInternalError("failed to delete template postings: " + err.Error())
    }
    for _, p := range postings {
        amt, _ := decimal.NewFromString(p.Amount)
        if _, err := q.CreateTemplatePosting(ctx(), generated.CreateTemplatePostingParams{
            TemplateID: id,
            AccountID:  p.AccountID,
            Currency:   p.Currency,
            Amount:     amt.String(),
        }); err != nil {
            return nil, model.NewInternalError("failed to create template posting: " + err.Error())
        }
    }

    // Replace tags
    if err := q.RemoveTemplateTags(ctx(), id); err != nil {
        return nil, model.NewInternalError("failed to remove template tags: " + err.Error())
    }
    for _, t := range tags {
        tagRow, err := q.CreateTag(ctx(), generated.CreateTagParams{Key: t.Key, Value: t.Value})
        if err != nil {
            return nil, model.NewInternalError("failed to create tag: " + err.Error())
        }
        if err := q.AddTemplateTag(ctx(), generated.AddTemplateTagParams{
            TemplateID: id,
            TagID:      tagRow.ID,
        }); err != nil {
            return nil, model.NewInternalError("failed to link tag: " + err.Error())
        }
    }

    if err := tx.Commit(); err != nil {
        return nil, model.NewInternalError("failed to commit: " + err.Error())
    }

    slog.Info("recurring template updated", "id", id)
    return s.GetTemplate(id)
}

// DeleteTemplate deletes a recurring template.
func (s *RecurringService) DeleteTemplate(id int64) error {
    _, err := s.queries.GetRecurringTemplate(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return model.NewNotFoundError("recurring template", id)
        }
        return model.NewInternalError("failed to get template: " + err.Error())
    }

    if err := s.queries.DeleteRecurringTemplate(ctx(), id); err != nil {
        return model.NewInternalError("failed to delete template: " + err.Error())
    }

    slog.Info("recurring template deleted", "id", id)
    return nil
}

// PreviewSchedule returns the next N dates for a given schedule configuration.
func (s *RecurringService) PreviewSchedule(scheduleType string, intervalN int, dayOfMonth *int, startDate string, endDate *string, count int) []string {
    start, err := time.Parse("2006-01-02", startDate)
    if err != nil {
        return nil
    }

    var end *time.Time
    if endDate != nil {
        e, err := time.Parse("2006-01-02", *endDate)
        if err == nil {
            end = &e
        }
    }

    var dates []string
    current := start

    for len(dates) < count {
        if end != nil && current.After(*end) {
            break
        }
        dates = append(dates, current.Format("2006-01-02"))
        current = advanceDate(current, scheduleType, intervalN, dayOfMonth)
    }

    return dates
}

// GenerateDueTransactions generates all overdue transactions for active templates up to the given date.
func (s *RecurringService) GenerateDueTransactions(asOfDate string) (int, error) {
    asOf, err := time.Parse("2006-01-02", asOfDate)
    if err != nil {
        return 0, model.NewValidationError("invalid date: " + asOfDate)
    }

    templates, err := s.queries.GetActiveTemplates(ctx())
    if err != nil {
        return 0, model.NewInternalError("failed to get active templates: " + err.Error())
    }

    totalGenerated := 0

    for _, tmpl := range templates {
        count, err := s.generateForTemplate(tmpl, asOf)
        if err != nil {
            slog.Error("failed to generate for template", "template_id", tmpl.ID, "error", err)
            continue
        }
        totalGenerated += count
    }

    if totalGenerated > 0 {
        slog.Info("recurring transactions generated", "count", totalGenerated)
    }
    return totalGenerated, nil
}

func (s *RecurringService) generateForTemplate(tmpl generated.RecurringTemplate, asOf time.Time) (int, error) {
    start, _ := time.Parse("2006-01-02", tmpl.StartDate)

    var end *time.Time
    if tmpl.EndDate.Valid {
        e, _ := time.Parse("2006-01-02", tmpl.EndDate.String)
        end = &e
    }

    // Determine where to start generating from
    var generateFrom time.Time
    if tmpl.LastGeneratedDate.Valid {
        lastGen, _ := time.Parse("2006-01-02", tmpl.LastGeneratedDate.String)
        // Advance one step past the last generated date
        var dayOfMonth *int
        if tmpl.DayOfMonth.Valid {
            d := int(tmpl.DayOfMonth.Int64)
            dayOfMonth = &d
        }
        generateFrom = advanceDate(lastGen, tmpl.ScheduleType, int(tmpl.IntervalN), dayOfMonth)
    } else {
        generateFrom = start
    }

    // Get template postings
    postingRows, err := s.queries.GetTemplatePostings(ctx(), tmpl.ID)
    if err != nil {
        return 0, fmt.Errorf("failed to get template postings: %w", err)
    }

    // Get template tags
    tagRows, err := s.queries.GetTemplateTags(ctx(), tmpl.ID)
    if err != nil {
        return 0, fmt.Errorf("failed to get template tags: %w", err)
    }

    postingInputs := make([]PostingInput, len(postingRows))
    for i, p := range postingRows {
        postingInputs[i] = PostingInput{
            AccountID: p.AccountID,
            Currency:  p.Currency,
            Amount:    p.Amount,
        }
    }

    tagInputs := make([]TagInput, len(tagRows))
    for i, t := range tagRows {
        tagInputs[i] = TagInput{Key: t.Key, Value: t.Value}
    }

    count := 0
    current := generateFrom
    var dayOfMonth *int
    if tmpl.DayOfMonth.Valid {
        d := int(tmpl.DayOfMonth.Int64)
        dayOfMonth = &d
    }

    for !current.After(asOf) {
        if end != nil && current.After(*end) {
            break
        }

        // Create the transaction
        _, err := s.txSvc.CreateTransaction(
            current.Format("2006-01-02"),
            tmpl.Description,
            tmpl.Status,
            postingInputs,
            tagInputs,
        )
        if err != nil {
            return count, fmt.Errorf("failed to create transaction for date %s: %w", current.Format("2006-01-02"), err)
        }

        // Update last_generated_date
        if err := s.queries.UpdateRecurringTemplateLastGenerated(ctx(), generated.UpdateRecurringTemplateLastGeneratedParams{
            LastGeneratedDate: sql.NullString{String: current.Format("2006-01-02"), Valid: true},
            UpdatedAt:         time.Now().Format(time.RFC3339),
            ID:                tmpl.ID,
        }); err != nil {
            return count, fmt.Errorf("failed to update last_generated_date: %w", err)
        }

        count++
        current = advanceDate(current, tmpl.ScheduleType, int(tmpl.IntervalN), dayOfMonth)
    }

    return count, nil
}

func advanceDate(current time.Time, scheduleType string, intervalN int, dayOfMonth *int) time.Time {
    switch scheduleType {
    case model.ScheduleEveryNDays:
        return current.AddDate(0, 0, intervalN)
    case model.ScheduleEveryNWeeks:
        return current.AddDate(0, 0, intervalN*7)
    case model.ScheduleEveryNMonths:
        return current.AddDate(0, intervalN, 0)
    case model.ScheduleMonthlyOnDay:
        next := current.AddDate(0, intervalN, 0)
        if dayOfMonth != nil {
            // Ensure the day of month is respected
            year, month, _ := next.Date()
            // Clamp to last day of month if needed
            lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
            day := *dayOfMonth
            if day > lastDay {
                day = lastDay
            }
            next = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
        }
        return next
    }
    return current.AddDate(0, 0, 1) // fallback
}

func validateSchedule(scheduleType string, intervalN int, dayOfMonth *int) error {
    switch scheduleType {
    case model.ScheduleEveryNDays, model.ScheduleEveryNWeeks, model.ScheduleEveryNMonths:
        // OK
    case model.ScheduleMonthlyOnDay:
        if dayOfMonth == nil {
            return model.NewValidationError("day_of_month is required for monthly_on_day schedule")
        }
        if *dayOfMonth < 1 || *dayOfMonth > 31 {
            return model.NewValidationError("day_of_month must be between 1 and 31")
        }
    default:
        return model.NewValidationError("invalid schedule_type: " + scheduleType)
    }
    if intervalN < 1 {
        return model.NewValidationError("interval_n must be at least 1")
    }
    return nil
}

func toModelTemplate(row generated.RecurringTemplate, postings []generated.RecurringTemplatePosting, tags []generated.Tag) *model.RecurringTemplate {
    var dayOfMonth *int
    if row.DayOfMonth.Valid {
        d := int(row.DayOfMonth.Int64)
        dayOfMonth = &d
    }
    var endDate *string
    if row.EndDate.Valid {
        endDate = &row.EndDate.String
    }
    var lastGenDate *string
    if row.LastGeneratedDate.Valid {
        lastGenDate = &row.LastGeneratedDate.String
    }

    modelPostings := make([]model.TemplatePosting, len(postings))
    for i, p := range postings {
        amt, _ := decimal.NewFromString(p.Amount)
        modelPostings[i] = model.TemplatePosting{
            ID:         p.ID,
            TemplateID: p.TemplateID,
            AccountID:  p.AccountID,
            Currency:   p.Currency,
            Amount:     model.StringDecimal{Decimal: amt},
        }
    }

    modelTags := make([]model.Tag, len(tags))
    for i, t := range tags {
        modelTags[i] = model.Tag{ID: t.ID, Key: t.Key, Value: t.Value}
    }

	return &model.RecurringTemplate{
		ID:                row.ID,
		Description:       row.Description,
		Status:            row.Status,
		ScheduleType:      row.ScheduleType,
		IntervalN:         int(row.IntervalN),
		DayOfMonth:        dayOfMonth,
		StartDate:         row.StartDate,
		EndDate:           endDate,
		LastGeneratedDate: lastGenDate,
		IsActive:          row.IsActive == 1,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
		Postings:          modelPostings,
		Tags:              modelTags,
	}
}

func toModelTemplateWithAccounts(row generated.RecurringTemplate, postings []generated.GetTemplatePostingsWithAccountsRow, tags []generated.Tag) *model.RecurringTemplate {
    var dayOfMonth *int
    if row.DayOfMonth.Valid {
        d := int(row.DayOfMonth.Int64)
        dayOfMonth = &d
    }
    var endDate *string
    if row.EndDate.Valid {
        endDate = &row.EndDate.String
    }
    var lastGenDate *string
    if row.LastGeneratedDate.Valid {
        lastGenDate = &row.LastGeneratedDate.String
    }

    modelPostings := make([]model.TemplatePosting, len(postings))
    for i, p := range postings {
        amt, _ := decimal.NewFromString(p.Amount)
        modelPostings[i] = model.TemplatePosting{
            ID:              p.ID,
            TemplateID:      p.TemplateID,
            AccountID:       p.AccountID,
            AccountFullName: p.AccountFullName,
            Currency:        p.Currency,
            Amount:          model.StringDecimal{Decimal: amt},
        }
    }

    modelTags := make([]model.Tag, len(tags))
    for i, t := range tags {
        modelTags[i] = model.Tag{ID: t.ID, Key: t.Key, Value: t.Value}
    }

    return &model.RecurringTemplate{
        ID:                row.ID,
        Description:       row.Description,
        Status:            row.Status,
        ScheduleType:      row.ScheduleType,
        IntervalN:         int(row.IntervalN),
        DayOfMonth:        dayOfMonth,
        StartDate:         row.StartDate,
        EndDate:           endDate,
        LastGeneratedDate: lastGenDate,
        IsActive:          row.IsActive == 1,
        CreatedAt:         row.CreatedAt,
        UpdatedAt:         row.UpdatedAt,
        Postings:          modelPostings,
        Tags:              modelTags,
    }
}
