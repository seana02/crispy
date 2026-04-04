package db

import (
	"fmt"
	"reflect"
	"testing"
)

// const transactions = []domain.Transaction{
// 	*domain.NewTransaction(-1, "A", getTime(2025, time.March, 9), domain.Cleared, nil, time.Now(), time.Now(), []*domain.Posting{}, []string{}),
// 	*domain.NewTransaction(-1, "B", getTime(2025, time.March, 9), domain.Cleared, nil, time.Now(), time.Now(), []*domain.Posting{}, []string{}),
// 	*domain.NewTransaction(-1, "C", getTime(2025, time.March, 9), domain.Cleared, nil, time.Now(), time.Now(), []*domain.Posting{}, []string{}),
// 	*domain.NewTransaction(-1, "D", getTime(2025, time.March, 9), domain.Cleared, nil, time.Now(), time.Now(), []*domain.Posting{}, []string{}),
// 	*domain.NewTransaction(-1, "E", getTime(2025, time.March, 9), domain.Cleared, nil, time.Now(), time.Now(), []*domain.Posting{}, []string{}),
// 	*domain.NewTransaction(-1, "F", getTime(2025, time.March, 9), domain.Cleared, nil, time.Now(), time.Now(), []*domain.Posting{}, []string{}),
// 	*domain.NewTransaction(-1, "G", getTime(2025, time.March, 9), domain.Cleared, nil, time.Now(), time.Now(), []*domain.Posting{}, []string{}),
// 	*domain.NewTransaction(-1, "H", getTime(2025, time.March, 9), domain.Cleared, nil, time.Now(), time.Now(), []*domain.Posting{}, []string{}),
// 	*domain.NewTransaction(-1, "I", getTime(2025, time.March, 9), domain.Cleared, nil, time.Now(), time.Now(), []*domain.Posting{}, []string{}),
// 	*domain.NewTransaction(-1, "J", getTime(2025, time.March, 9), domain.Cleared, nil, time.Now(), time.Now(), []*domain.Posting{}, []string{}),
// }

// Ensure constructor and methods build struct properly
func Test_TransactionBuilderStruct(t *testing.T) {
	s, err := InitSQLite(":memory:")
	if err != nil {
		t.Fatalf("Error initializing in-memory SQLite!: %+v", err)
	}
	if s == nil {
		t.Fatalf("Received nil sqlite object: %+v", err)
	}

	columns := []TransactionColumn{
		TransactionColumn_Id,
		TransactionColumn_Description,
		TransactionColumn_Date,
		TransactionColumn_Status,
		TransactionColumn_ReferenceId,
		TransactionColumn_DateCreated,
		TransactionColumn_DateUpdated,
		TransactionColumn_All,
	}
	// SELECT [column] FROM "Transaction"
	for _, c := range columns {
		t.Run(fmt.Sprintf("SELECT %s", c), func(t *testing.T) {
			got := s.TransactionBuilder(StatementType_Select)
			want := &SQLite_TransactionBuilder{
				handle:        s,
				statementType: StatementType_Select,
				columns:       make(map[TransactionColumn]struct{}),
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Expected %+v got %+v", want, got)
			}
		})
	}

	addColumn := [][]TransactionColumn{
		[]TransactionColumn{TransactionColumn_Id, TransactionColumn_Description},
		[]TransactionColumn{TransactionColumn_Id, TransactionColumn_Id},
		[]TransactionColumn{TransactionColumn_Date, TransactionColumn_Id, TransactionColumn_Description, TransactionColumn_Date},
	}
	for i, list := range addColumn {
		t.Run(fmt.Sprintf("Build columns %d", i), func(t *testing.T) {
			desiredMap := make(map[TransactionColumn]struct{})
			got := s.TransactionBuilder(StatementType_Select)
			for _, c := range list {
				got = got.AddColumn(c)
				desiredMap[c] = struct{}{}
			}
			want := &SQLite_TransactionBuilder{
				handle:        s,
				statementType: StatementType_Select,
				columns:       desiredMap,
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Expected %+v got %+v", want, got)
			}
		})
	}

	whereList := [][]Where{
		[]Where{{columns[0], Equal, "2"}},
		[]Where{{columns[1], Like, "%Description%"},
			{columns[2], GreaterThan, "2026-01-02"}},
	}
	for i, list := range whereList {
		t.Run(fmt.Sprintf("Build where %d", i), func(t *testing.T) {
			got := s.TransactionBuilder(StatementType_Select)
			for _, w := range list {
				got = got.AddCondition(w.Column, w.Op, w.Value)
			}
			want := &SQLite_TransactionBuilder{
				handle:        s,
				statementType: StatementType_Select,
				columns:       make(map[TransactionColumn]struct{}),
				where:         list,
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("Expected %+v got %+v", want, got)
			}
		})
	}
}
