package service

import (
    "database/sql"
    "log/slog"
    "strings"
    "time"

    "github.com/crispy-2/backend/internal/dal/generated"
    "github.com/crispy-2/backend/internal/model"
)

// AccountService handles account business logic.
type AccountService struct {
    db      *sql.DB
    queries *generated.Queries
}

// NewAccountService creates a new AccountService.
func NewAccountService(db *sql.DB) *AccountService {
    return &AccountService{
        db:      db,
        queries: generated.New(db),
    }
}

// CreateAccount creates a new account.
func (s *AccountService) CreateAccount(name string, parentID *int64, accountType string) (*model.Account, error) {
    if err := validateAccountName(name); err != nil {
        return nil, err
    }

    var fullName string
    var parentIDSQL sql.NullInt64
    var accountTypeSQL sql.NullString

    if parentID != nil {
        parentIDSQL = sql.NullInt64{Int64: *parentID, Valid: true}

        // Fetch parent to build full_name
        parent, err := s.queries.GetAccount(ctx(), *parentID)
        if err != nil {
            if err == sql.ErrNoRows {
                return nil, model.NewNotFoundError("account", *parentID)
            }
            return nil, model.NewInternalError("failed to fetch parent account: " + err.Error())
        }
        fullName = parent.FullName + ":" + name

        // Child accounts must not specify a type; it is copied from the parent
        if accountType != "" {
            return nil, model.NewValidationError("child accounts inherit their type from the root ancestor; do not specify account_type")
        }
        accountTypeSQL = parent.AccountType
    } else {
        fullName = name

        // Root accounts must have a type
        if accountType == "" {
            return nil, model.NewValidationError("root accounts must have an account_type")
        }
        if !model.ValidAccountType(accountType) {
            return nil, model.NewValidationError("invalid account_type: " + accountType)
        }
        accountTypeSQL = sql.NullString{String: accountType, Valid: true}
    }

    row, err := s.queries.CreateAccount(ctx(), generated.CreateAccountParams{
        Name:        name,
        FullName:    fullName,
        ParentID:    parentIDSQL,
        AccountType: accountTypeSQL,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to create account: " + err.Error())
    }

    slog.Info("account created", "id", row.ID, "full_name", fullName)
    return toModelAccount(row), nil
}

// CreateAccountByPath creates an account from a full ":" separated path like
// "Expenses:Food:Restaurants". It creates any missing intermediate accounts.
// accountType is used only if the root-level account must be created; it is
// auto-suggested from the root segment name if empty.
func (s *AccountService) CreateAccountByPath(fullName string, accountType string) (*model.Account, error) {
    if fullName == "" {
        return nil, model.NewValidationError("account name must not be empty")
    }

    segments := strings.Split(fullName, ":")
    for _, seg := range segments {
        if strings.TrimSpace(seg) == "" {
            return nil, model.NewValidationError("account path contains an empty segment")
        }
    }

    // If there's only one segment, this is a root account.
    if len(segments) == 1 {
        // Prefer suggestion over provided type, fall back to provided type.
        suggested := string(model.SuggestAccountType(segments[0]))
        if suggested != "" {
            accountType = suggested
        } else if accountType == "" {
            return nil, model.NewValidationError("root accounts must have an account_type")
        }
        return s.CreateAccount(segments[0], nil, accountType)
    }

    // Walk the hierarchy, creating missing segments.
    var parentID *int64
    var currentFullName string
    var leaf *model.Account

    for depth, seg := range segments {
        if depth == 0 {
            currentFullName = seg
        } else {
            currentFullName = currentFullName + ":" + seg
        }

        // Check if this level already exists
        existing, err := s.SearchAccounts(currentFullName)
        if err == nil {
            for _, a := range existing {
                if a.FullName == currentFullName {
                    id := a.ID
                    parentID = &id
                    leaf = &a
                    goto nextSegment
                }
            }
        }

        // Doesn't exist — create it.
        {
            segAccountType := ""
            if parentID == nil {
                // Root segment needs a type. Prefer suggestion, fall back to provided default.
                suggested := string(model.SuggestAccountType(seg))
                if suggested != "" {
                    segAccountType = suggested
                } else if accountType != "" {
                    segAccountType = accountType
                } else {
                    return nil, model.NewValidationError("could not determine account type for root '" + seg + "'; specify account_type")
                }
            }
            newAcct, err := s.CreateAccount(seg, parentID, segAccountType)
            if err != nil {
                return nil, model.NewValidationError("could not create account '" + currentFullName + "': " + err.Error())
            }
            id := newAcct.ID
            parentID = &id
            leaf = newAcct
        }

    nextSegment:
    }

    return leaf, nil
}

// GetAccount retrieves a single account by ID.
func (s *AccountService) GetAccount(id int64) (*model.Account, error) {
    row, err := s.queries.GetAccount(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, model.NewNotFoundError("account", id)
        }
        return nil, model.NewInternalError("failed to get account: " + err.Error())
    }
    return toModelAccount(row), nil
}

// ListAccounts returns a paginated list of accounts with optional filters.
func (s *AccountService) ListAccounts(nameFilter, typeFilter string, pg model.Pagination) (*model.PaginatedResult[model.Account], error) {
    pg = normalizePagination(pg)
    offset := int64((pg.Page - 1) * pg.PageSize)
    limit := int64(pg.PageSize)

    var rows []generated.Account
    var total int64
    var err error

    hasName := nameFilter != ""
    hasType := typeFilter != ""

    switch {
    case hasName && hasType:
        namePattern := "%" + nameFilter + "%"
        rows, err = s.queries.ListAccountsByTypeAndName(ctx(), generated.ListAccountsByTypeAndNameParams{
            AccountType: sql.NullString{String: typeFilter, Valid: true},
            FullName:    namePattern,
            Limit:       limit,
            Offset:      offset,
        })
        if err == nil {
            total, err = s.queries.CountAccountsByTypeAndName(ctx(), generated.CountAccountsByTypeAndNameParams{
                AccountType: sql.NullString{String: typeFilter, Valid: true},
                FullName:    namePattern,
            })
        }
    case hasName:
        namePattern := "%" + nameFilter + "%"
        rows, err = s.queries.ListAccountsByName(ctx(), generated.ListAccountsByNameParams{
            FullName: namePattern,
            Limit:    limit,
            Offset:   offset,
        })
        if err == nil {
            total, err = s.queries.CountAccountsByName(ctx(), namePattern)
        }
    case hasType:
        rows, err = s.queries.ListAccountsByType(ctx(), generated.ListAccountsByTypeParams{
            AccountType: sql.NullString{String: typeFilter, Valid: true},
            Limit:       limit,
            Offset:      offset,
        })
        if err == nil {
            total, err = s.queries.CountAccountsByType(ctx(), sql.NullString{String: typeFilter, Valid: true})
        }
    default:
        rows, err = s.queries.ListAccounts(ctx(), generated.ListAccountsParams{
            Limit:  limit,
            Offset: offset,
        })
        if err == nil {
            total, err = s.queries.CountAccounts(ctx())
        }
    }

    if err != nil {
        return nil, model.NewInternalError("failed to list accounts: " + err.Error())
    }

    items := make([]model.Account, len(rows))
    for i, r := range rows {
        items[i] = *toModelAccount(r)
    }

    totalPages := int(total) / pg.PageSize
    if int(total)%pg.PageSize != 0 {
        totalPages++
    }

    return &model.PaginatedResult[model.Account]{
        Items:      items,
        Total:      int(total),
        Page:       pg.Page,
        PageSize:   pg.PageSize,
        TotalPages: totalPages,
    }, nil
}

// UpdateAccount updates an account's name, parent, and/or type.
func (s *AccountService) UpdateAccount(id int64, name *string, parentID **int64, accountType string) (*model.Account, error) {
    existing, err := s.queries.GetAccount(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, model.NewNotFoundError("account", id)
        }
        return nil, model.NewInternalError("failed to get account: " + err.Error())
    }

    newName := existing.Name
    if name != nil {
        if err := validateAccountName(*name); err != nil {
            return nil, err
        }
        newName = *name
    }

    newParentID := existing.ParentID
    if parentID != nil {
        if *parentID != nil {
            newParentID = sql.NullInt64{Int64: **parentID, Valid: true}
        } else {
            newParentID = sql.NullInt64{}
        }
    }

    newAccountType := existing.AccountType
    if accountType != "" {
        if !model.ValidAccountType(accountType) {
            return nil, model.NewValidationError("invalid account_type: " + accountType)
        }
        newAccountType = sql.NullString{String: accountType, Valid: true}
    }

    // Compute new full_name
    var newFullName string
    if newParentID.Valid {
        parent, err := s.queries.GetAccount(ctx(), newParentID.Int64)
        if err != nil {
            return nil, model.NewInternalError("failed to get parent: " + err.Error())
        }
        newFullName = parent.FullName + ":" + newName
    } else {
        newFullName = newName
    }

    row, err := s.queries.UpdateAccount(ctx(), generated.UpdateAccountParams{
        Name:        newName,
        FullName:    newFullName,
        ParentID:    newParentID,
        AccountType: newAccountType,
        UpdatedAt:   time.Now().Format(time.RFC3339),
        ID:          id,
    })
    if err != nil {
        return nil, model.NewInternalError("failed to update account: " + err.Error())
    }

    // Update all descendants' full_name if the name or parent changed
    if newFullName != existing.FullName {
        if err := s.updateDescendantFullNames(existing.FullName, newFullName); err != nil {
            return nil, err
        }
    }

    slog.Info("account updated", "id", id, "full_name", newFullName)
    return toModelAccount(row), nil
}

// DeleteAccount deletes an account if it has no children or postings.
func (s *AccountService) DeleteAccount(id int64) error {
    // First verify account exists
    _, err := s.queries.GetAccount(ctx(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            return model.NewNotFoundError("account", id)
        }
        return model.NewInternalError("failed to get account: " + err.Error())
    }

    // Check for children
    children, err := s.queries.GetChildAccounts(ctx(), sql.NullInt64{Int64: id, Valid: true})
    if err != nil {
        return model.NewInternalError("failed to check children: " + err.Error())
    }
    if len(children) > 0 {
        return model.NewConflictError("cannot delete account: it has child accounts")
    }

    // Check for postings
    hasPostings, err := s.queries.HasPostingsForAccount(ctx(), id)
    if err != nil {
        return model.NewInternalError("failed to check postings: " + err.Error())
    }
    if hasPostings {
        return model.NewConflictError("cannot delete account: it has postings referencing it")
    }

    if err := s.queries.DeleteAccount(ctx(), id); err != nil {
        return model.NewInternalError("failed to delete account: " + err.Error())
    }

    slog.Info("account deleted", "id", id)
    return nil
}

// SearchAccounts performs a fuzzy search on account full_names.
func (s *AccountService) SearchAccounts(query string) ([]model.Account, error) {
    pattern := "%" + strings.ToLower(query) + "%"
    rows, err := s.queries.SearchAccounts(ctx(), pattern)
    if err != nil {
        return nil, model.NewInternalError("failed to search accounts: " + err.Error())
    }

    accounts := make([]model.Account, len(rows))
    for i, r := range rows {
        accounts[i] = *toModelAccount(r)
    }
    return accounts, nil
}

// updateDescendantFullNames updates full_name for all descendants when an ancestor is renamed.
func (s *AccountService) updateDescendantFullNames(oldPrefix, newPrefix string) error {
    descendants, err := s.queries.GetDescendantAccounts(ctx(), oldPrefix+":%" )
    if err != nil {
        return model.NewInternalError("failed to get descendants: " + err.Error())
    }
    for _, d := range descendants {
        newFullName := newPrefix + strings.TrimPrefix(d.FullName, oldPrefix)
    _, err := s.queries.UpdateAccount(ctx(), generated.UpdateAccountParams{
        Name:        d.Name,
        FullName:    newFullName,
        ParentID:    d.ParentID,
        AccountType: d.AccountType,
        UpdatedAt:   time.Now().Format(time.RFC3339),
        ID:          d.ID,
    })
        if err != nil {
            return model.NewInternalError("failed to update descendant full_name: " + err.Error())
        }
    }
    return nil
}

func validateAccountName(name string) error {
    if name == "" {
        return model.NewValidationError("account name must not be empty")
    }
    if strings.Contains(name, ":") {
        return model.NewValidationError("account name must not contain ':'")
    }
    return nil
}

func toModelAccount(row generated.Account) *model.Account {
	var parentID *int64
	if row.ParentID.Valid {
		parentID = &row.ParentID.Int64
	}
	return &model.Account{
		ID:          row.ID,
		Name:        row.Name,
		FullName:    row.FullName,
		ParentID:    parentID,
		AccountType: row.AccountType.String,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
