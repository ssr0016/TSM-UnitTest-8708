package schedulerimpl

import (
	"context"
	"database/sql"
	"fmt"
	dbtesting "main/pkg/infra/storage/db/testing"
	"main/pkg/infra/storage/postgres"
	"main/pkg/tsm/scheduler"
	"main/pkg/util/pointer"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		name           string
		expectedResult int64
		entity         *scheduler.Scheduler
		expectedError  error
	}{
		{
			name:           "create scheduler Success",
			expectedResult: 1,
			entity:         &scheduler.Scheduler{},
			expectedError:  nil,
		},
		{
			name:           "create scheduler Error",
			expectedResult: 0,
			entity:         &scheduler.Scheduler{},
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

			_, err := store.create(context.Background(), tc.entity)
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestCreateAssignee(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult interface{}
		entity         *scheduler.SchedulerAssignee
	}{
		{
			name:           "create scheduler assignee Success",
			expectedError:  nil,
			expectedResult: 1,
			entity:         &scheduler.SchedulerAssignee{},
		},
		{
			name:           "create assigment log Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: 0,
			entity:         &scheduler.SchedulerAssignee{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)

			err := store.createAssignee(context.Background(), tc.entity)
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		result         []*scheduler.Scheduler
		query          *scheduler.SearchSchedulerQuery
		expectedResult *scheduler.SearchSchedulerQueryResult
	}{
		{
			name:          "Success",
			expectedError: nil,
			result:        []*scheduler.Scheduler{},
			expectedResult: &scheduler.SearchSchedulerQueryResult{
				Schedulers: []*scheduler.Scheduler{},
			},
			query: &scheduler.SearchSchedulerQuery{
				Name:      "name",
				Currency:  "currency",
				Priority:  1,
				Status:    1,
				Assignees: []int64{1},
				DateFrom:  &time.Time{},
				DateTo:    &time.Time{},
				PerPage:   10,
			},
		},
		{
			name:           "Error",
			query:          &scheduler.SearchSchedulerQuery{},
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

func TestGetSchedulerByID(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult *scheduler.SchedulerDTO
		schedulerID    int64
	}{
		{
			name:           "get assigment by id Success",
			expectedError:  nil,
			expectedResult: &scheduler.SchedulerDTO{},
			schedulerID:    1,
		},
		{
			name:           "get assigment by id Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			schedulerID:    1,
		},
		{
			name:           "get assigment by id Error - Not Found",
			expectedError:  sql.ErrNoRows,
			expectedResult: nil,
			schedulerID:    1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)

			result, err := store.GetSchedulerByID(context.Background(), tc.schedulerID)
			if err == nil && result != nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGetSchedulerAssignByID(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult []*scheduler.SchedulerAssignee
		assignID       int64
	}{
		{
			name:           "Success",
			expectedError:  nil,
			expectedResult: []*scheduler.SchedulerAssignee{},
			assignID:       1,
		},
		{
			name:           "Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			assignID:       1,
		},
		{
			name:           "Not Found",
			expectedError:  sql.ErrNoRows,
			expectedResult: nil,
			assignID:       1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)

			result, err := store.GetSchedulerAssignByID(context.Background(), tc.assignID)
			if err == nil && result != nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestGetSchedulerUserIDsByID(t *testing.T) {
	testCases := []struct {
		name           string
		userIDsByID    []int64
		expectedResult []*int64
		expectedError  error
	}{
		{
			name:           "Success",
			userIDsByID:    []int64{1},
			expectedResult: []*int64{pointer.Int64Ptr(1)},
			expectedError:  nil,
		},
		{
			name:           "Error",
			userIDsByID:    []int64{0},
			expectedResult: nil,
			expectedError:  fmt.Errorf("error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedResult: tc.expectedResult,
				ExpectedError:  tc.expectedError,
			}

			store := newStore(fakeDB)

			result, err := store.GetSchedulerUserIDsByID(context.Background(), tc.userIDsByID[0])
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestUpdateScheduleStatus(t *testing.T) {
	testCases := []struct {
		name           string
		expectedError  error
		expectedResult interface{}
		entity         *scheduler.SchedulerDTO
	}{
		{
			name:           "update scheduler status Success",
			expectedError:  nil,
			expectedResult: nil,
			entity:         &scheduler.SchedulerDTO{},
		},
		{
			name:           "update scheduler status Error",
			expectedError:  fmt.Errorf("error"),
			expectedResult: nil,
			entity:         &scheduler.SchedulerDTO{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError:  tc.expectedError,
				ExpectedResult: tc.expectedResult,
			}

			store := newStore(fakeDB)

			err := store.updateScheduleStatus(context.Background(), tc.entity)
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

		})
	}
}

func TestUnassignAssignee(t *testing.T) {
	testCases := []struct {
		name          string
		id            int64
		expectedError error
	}{
		{
			name:          "unassign assignee Success",
			id:            1,
			expectedError: nil,
		},
		{
			name:          " unassign assignee Error",
			id:            1,
			expectedError: fmt.Errorf("error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError: tc.expectedError,
			}

			store := newStore(fakeDB)

			err := store.unassignAssignee(context.Background(), tc.id)
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

		})
	}
}

func TestUpdateScheduler(t *testing.T) {
	testCases := []struct {
		name          string
		expectedError error
		entity        *scheduler.Scheduler
	}{
		{
			name:          "Success",
			expectedError: nil,
			entity:        &scheduler.Scheduler{},
		},
		{
			name:          "Error",
			expectedError: fmt.Errorf("error"),
			entity:        &scheduler.Scheduler{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError: tc.expectedError,
			}

			store := newStore(fakeDB)

			err := store.updateScheduler(context.Background(), tc.entity)
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestDeleteBySchedulerIDAndUserID(t *testing.T) {
	testCases := []struct {
		name          string
		id            int64
		userID        int64
		expectedError error
	}{
		{
			name:          "Success",
			id:            1,
			userID:        1,
			expectedError: nil,
		},
		{
			name:          "Error",
			id:            1,
			userID:        1,
			expectedError: fmt.Errorf("error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeDB := &dbtesting.FakeSqlxdb{
				ExpectedError: tc.expectedError,
			}

			store := newStore(fakeDB)

			err := store.deleteBySchedulerIDAndUserID(context.Background(), tc.id, tc.userID)
			if err == nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpected error: %v", err)
			}

		})
	}
}

func TestGetSchedulerList(t *testing.T) {
	testCases := []struct {
		name           string
		id             int64
		expectedResult []*scheduler.SchedulerDTO
		expectedError  error
	}{
		{
			name:           "Success",
			expectedError:  nil,
			expectedResult: []*scheduler.SchedulerDTO{},
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

			result, err := store.getSchedulerList(context.Background())
			if err == nil && result != nil && tc.expectedError != nil {
				t.Fatalf("expected error %q, but got none", tc.expectedError)
			} else if err != nil && tc.expectedError == nil {
				t.Fatalf("unexpectred error: %v", err)
			}

			require.Equal(t, tc.expectedResult, result)
		})
	}
}

func Test_store_GetSchedulerUserIDsByID(t *testing.T) {
	type fields struct {
		db     postgres.DB
		logger *zap.Logger
	}
	type args struct {
		ctx         context.Context
		schedulerID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &store{
				db:     tt.fields.db,
				logger: tt.fields.logger,
			}
			got, err := s.GetSchedulerUserIDsByID(tt.args.ctx, tt.args.schedulerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("store.GetSchedulerUserIDsByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("store.GetSchedulerUserIDsByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
