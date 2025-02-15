package ports

import "context"

type UserRepositoryInterface interface {
	GetUserURLS(ctx context.Context, id string)
}
