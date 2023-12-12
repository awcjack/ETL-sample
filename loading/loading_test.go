package loading_test

import (
	"context"
	"testing"
	"time"

	"github.com/awcjack/ETL-sample/loading"
	"github.com/awcjack/ETL-sample/transformation"
)

func TestSaveData(t *testing.T) {
	type testcase struct {
		testcase             string
		dataCount            int
		bulkInsert           bool
		bulkInsertSize       int
		bulkInsertInterval   int
		expectedRepoAddUser  int
		expectedRepoAddUsers int
	}

	testcases := []testcase{
		{
			testcase:             "Single Insert",
			dataCount:            5,
			bulkInsert:           false,
			bulkInsertSize:       10,
			bulkInsertInterval:   10,
			expectedRepoAddUser:  5,
			expectedRepoAddUsers: 0,
		},
		{
			testcase:             "Bulk Insert",
			dataCount:            5,
			bulkInsert:           true,
			bulkInsertSize:       5,
			bulkInsertInterval:   1000, // hard to test code based on time
			expectedRepoAddUser:  0,
			expectedRepoAddUsers: 1,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.testcase, func(t *testing.T) {
			dep := newLoadingDependencies()
			dataChan := make(chan transformation.TransformedData, tc.dataCount)
			for i := 0; i < tc.dataCount; i++ {
				dataChan <- transformation.TransformedData{}
			}

			ctx, cancel := context.WithCancel(context.Background())

			go func() {
				// force end goroutine after 3 secnd
				time.Sleep(3 * time.Second)
				cancel()
			}()

			loading.SaveData(ctx, dep.repo, dataChan, tc.bulkInsert, tc.bulkInsertSize, tc.bulkInsertInterval)
			if dep.repo.CalledAddUser != tc.expectedRepoAddUser {
				t.Errorf("expected called add user %v, but got %v", tc.expectedRepoAddUser, dep.repo.CalledAddUser)
			}
			if dep.repo.CalledAddUsers != tc.expectedRepoAddUsers {
				t.Errorf("expected called add users %v, but got %v", tc.expectedRepoAddUsers, dep.repo.CalledAddUsers)
			}
		})
	}
}

type loadingDependencies struct {
	repo *repoMock
}

func newLoadingDependencies() loadingDependencies {
	repo := &repoMock{}

	return loadingDependencies{
		repo: repo,
	}
}

type repoMock struct {
	CalledAddUser  int
	CalledAddUsers int
}

func (r *repoMock) AddUser(ctx context.Context, user transformation.TransformedData) error {
	r.CalledAddUser++
	return nil
}

func (r *repoMock) AddUsers(ctx context.Context, user []transformation.TransformedData) error {
	r.CalledAddUsers++
	return nil
}
