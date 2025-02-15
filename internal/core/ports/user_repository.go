package ports

import "context"

type UserRepository interface {
	GetUserURLS(ctx context.Context, id string)
}
