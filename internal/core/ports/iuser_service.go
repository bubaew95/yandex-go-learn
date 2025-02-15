package ports

import "context"

type UserServiceInterface interface {
	GetUserURLS(ctx context.Context, id string)
}
