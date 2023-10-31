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

func TestCreateAssignment(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult int64
		entity         *assignment.Assignment
	}{
		{
			name:           "create assignment Success",
			expectedError:  nil,
			expectedResult: 1,
			entity:         &assignment.Assignment{},
		},
		{
			name:           "create assignment Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: 0,
			entity:         &assignment.Assignment{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)
			result, err := store.create(context.Background(), tc.entity)

			if err == nil && tc.expectedError != nil {
				t.Fatalf("expeted error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpacted error %v", err)
			}

			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestCreateAssignmentLog(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResukt interface{}
		entity         *assignment.AssignmentLog
	}{
		{
			name:           "create assignment log Success",
			expectedError:  nil,
			expectedResukt: 1,
			entity:         &assignment.AssignmentLog{},
		},
		{
			name:           "create assignment log Error",
			expectedError:  fmt.Errorf("error"),
			expectedResukt: 0,
			entity:         &assignment.AssignmentLog{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResukt,
			}

			store := newStore(fakeDB)
			err := store.createAssignmentLog(context.Background(), tc.entity)

			if err == nil && tc.expectedError != nil {
				t.Fatalf("expeted error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpacted error %v", err)
			}

		})
	}
}

func TestGetAssignmentByMemberID(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult *assignment.Assignment
		memberID       int64
	}{
		{
			name:           "Success",
			expectedError:  nil,
			expectedResult: &assignment.Assignment{},
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
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}
			require.Equal(t, tc.expectedResult, result)
		})
	}

}

func TestGetAssignmentByID(t *testing.T) {
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
			name:           "get batch adjustment by id Error- Not Found",
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

func TestSearchAssignment(t *testing.T) {

	testCases := []struct {
		name           string
		expectedError  error
		expectedResult *assignment.SearchAssignmentQueryResult
		query          *assignment.SearchAssignmentQuery
	}{
		{
			name:           "search assigment Success",
			expectedError:  nil,
			expectedResult: &assignment.SearchAssignmentQueryResult{},
			query: &assignment.SearchAssignmentQuery{
				MemberID:  1,
				Assignees: []int64{1},
				DateFrom:  &time.Time{},
				DateTo:    &time.Time{},
				PerPage:   10,
			},
		},
		{
			name:           "search assignment Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			query:          &assignment.SearchAssignmentQuery{},
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

func TestGetAssignmentCount(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult *assignment.Assignment
		assignmentID   int64
	}{
		{
			name:           "Success",
			expectedError:  nil,
			expectedResult: &assignment.Assignment{},
			assignmentID:   1,
		},
		{
			name:           "Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			assignmentID:   1,
		},
		{
			name:           "Not Found",
			expectedError:  fmt.Errorf("error"),
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
			result, err := store.getByMemberID(context.Background(), tc.assignmentID)
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
			name:           "get assignment by id Success",
			expectedError:  nil,
			expectedResult: []int64{1},
			assignmentID:   []int64{1},
		},
		{
			name:           "get assignment by id Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			assignmentID:   []int64{1},
		},
		{
			name:           "get batch adjustment by id Error- Not Found",
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

func TestUpdateAssignment(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult interface{}
		entity         *assignment.Assignment
	}{
		{
			name:           "Success",
			expectedError:  nil,
			expectedResult: nil,
			entity:         &assignment.Assignment{},
		},
		{
			name:           "Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			entity:         &assignment.Assignment{},
		},
		{
			name:           "Not Found",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			entity:         &assignment.Assignment{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
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

func getAssignmentLog(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult *assignment.AssignmentLog
		entity         int64
	}{
		{
			name:           "Success",
			expectedError:  nil,
			expectedResult: nil,
			entity:         1,
		},
		{
			name:           "Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			entity:         1,
		},
		{
			name:           "Not Found",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			entity:         1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)
			result, err := store.getAssignmentLog(context.Background(), tc.entity)

			if err == nil && result != nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpectred error: %v", err)
			}

			require.Equal(t, tc.expectedResult, result)
		})
	}

}

func Test_getAssignmentLog(t *testing.T) {
	type args struct {
		t *testing.T
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAssignmentLog(tt.args.t)
		})
	}
}
