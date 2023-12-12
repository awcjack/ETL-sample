package loading

import (
	"context"

	"github.com/awcjack/ETL-sample/transformation"
)

type repository interface {
	AddUser(ctx context.Context, user transformation.TransformedData) error
	AddUsers(ctx context.Context, users []transformation.TransformedData) error
}
