package loading

import (
	"context"
	"time"

	"github.com/awcjack/ETL-sample/transformation"
)

// Saving data to specific repo
// allow bulk insert or single insert
// flush based on if slice didn't filled
func SaveData(ctx context.Context, repo repository, dataPipeline <-chan transformation.TransformedData, bulkInsert bool, bulkInsertSize int, bulkInsertInterval int) error {
	var err error
	if bulkInsert {
		timer := time.NewTimer(time.Duration(bulkInsertInterval) * time.Second)
		users := make([]transformation.TransformedData, 0, bulkInsertSize)
		for {
			select {
			case <-ctx.Done():
				return nil
			case <-timer.C:
				if len(users) != 0 {
					err = repo.AddUsers(ctx, users)
					users = users[:0]
					if err != nil {
						return err
					}
				}
				timer.Reset(time.Duration(bulkInsertInterval) * time.Second)
			case data := <-dataPipeline:
				users = append(users, data)
				if len(users) >= bulkInsertSize {
					err = repo.AddUsers(ctx, users)
					users = users[:0]
					if err != nil {
						return err
					}
					timer.Reset(time.Duration(bulkInsertInterval) * time.Second)
				}
			}
		}
	} else {
		for {
			select {
			case <-ctx.Done():
				return nil
			case data := <-dataPipeline:
				err = repo.AddUser(ctx, data)
				if err != nil {
					return err
				}
			}
		}
	}
}
