package assignmentimpl

import (
	"context"
	"database/sql"
	"fmt"
	dbtesting "main/pkg/infra/storage/db/testing"
	"main/pkg/tsm/assignment"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult int64
		entity         *assignment.Assignment
		expectedError  error
	}{
		{
			name:           "create assignment Success",
			expectedResult: 1,
			entity:         &assignment.Assignment{},
			expectedError:  nil,
		},
		{
			name:           "create assignment Error",
			expectedResult: 0,
			entity:         &assignment.Assignment{},
			expectedError:  fmt.Errorf("error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError: tc.expectedError,
			}

			store := newStore(fakeDB)

			_, err := store.create(context.Background(), tc.entity)
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error %v", err)
			}
		})
	}
}

func TestCreateAssignmentLog(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult interface{}
		entity         *assignment.AssignmentLog
	}{
		{
			name:           "create assignment log Success",
			expectedError:  nil,
			expectedResult: 1,
			entity:         &assignment.AssignmentLog{},
		},
		{
			name:           "create assignment log Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: 0,
			entity:         &assignment.AssignmentLog{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)

			err := store.createAssignmentLog(context.Background(), tc.entity)
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpacted error %v", err)
			}
		})
	}
}

func TestGetByMemberID(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult *assignment.AssignmentDTO
		memberID       int64
	}{
		{
			name:           "Success",
			expectedError:  nil,
			expectedResult: &assignment.AssignmentDTO{},
			memberID:       1,
		},
		{
			name:           "Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			memberID:       1,
		},
		{
			name:           "Not Found",
			expectedError:  sql.ErrNoRows,
			expectedResult: nil,
			memberID:       1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)

			result, err := store.getByMemberID(context.Background(), tc.memberID)
			if err == nil && result != nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

			require.Equal(t, tc.expectedResult, result)
		})
	}

}

func TestGetByID(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult *assignment.AssignmentDTO
		assignmentID   int64
	}{
		{
			name:           "get assignment by id Success",
			expectedError:  nil,
			expectedResult: &assignment.AssignmentDTO{},
			assignmentID:   1,
		},
		{
			name:           "get assignment by id Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			assignmentID:   1,
		},
		{
			name:           "get assigment by id Error- Not Found",
			expectedError:  sql.ErrNoRows,
			expectedResult: nil,
			assignmentID:   1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)

			result, err := store.getByID(context.Background(), tc.assignmentID)
			if err == nil && result != nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestSearch(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		result         []*assignment.Assignment
		expectedResult *assignment.SearchAssignmentQueryResult
		query          *assignment.SearchAssignmentQuery
	}{
		{
			name:          "Success",
			expectedError: nil,
			result:        []*assignment.Assignment{},
			expectedResult: &assignment.SearchAssignmentQueryResult{
				Assignments: []*assignment.Assignment{},
			},
			query: &assignment.SearchAssignmentQuery{
				MemberID:  1,
				Assignees: []int64{1},
				DateFrom:  &time.Time{},
				DateTo:    &time.Time{},
				PerPage:   10,
			},
		},
		{
			name:           "Error",
			query:          &assignment.SearchAssignmentQuery{},
			expectedResult: nil,
			expectedError:  fmt.Errorf("error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)

			result, err := store.search(context.Background(), tc.query)
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGetByAssigneesID(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult []int64
		assignmentID   []int64
	}{
		{
			name:           "Success",
			expectedError:  nil,
			expectedResult: []int64{1},
			assignmentID:   []int64{1},
		},
		{
			name:           "Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			assignmentID:   []int64{1},
		},
		{
			name:           "Not Found",
			expectedError:  sql.ErrNoRows,
			expectedResult: nil,
			assignmentID:   []int64{1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)

			result, err := store.getByAssigneesID(context.Background(), tc.assignmentID)
			if err == nil && result != nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

			require.Equal(t, tc.expectedResult, result)
		})
	}

}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name          string
		expectedError error
		entity        *assignment.Assignment
	}{
		{
			name:          "Success",
			expectedError: nil,
			entity:        &assignment.Assignment{},
		},
		{
			name:          "Error",
			expectedError: fmt.Errorf("error"),
			entity:        &assignment.Assignment{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError: tc.expectedError,
			}

			store := newStore(fakeDB)

			err := store.update(context.Background(), tc.entity)
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetAssignmentLog(t *testing.T) {
	testCases := []struct {
		name           string
		id             int64
		expectedResult []*assignment.AssignmentLog
		expectedError  error
	}{
		{
			name:           "Success",
			expectedError:  nil,
			expectedResult: []*assignment.AssignmentLog{},
			id:             1,
		},
		{
			name:           "Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			id:             1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)

			result, err := store.getAssignmentLog(context.Background(), tc.id)
			if err == nil && result != nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpectred error: %v", err)
			}

			require.Equal(t, tc.expectedResult, result)
		})
	}
}
