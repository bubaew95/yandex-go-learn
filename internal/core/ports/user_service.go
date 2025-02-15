package ports

import "context"

type UserService interface {
	GetUserURLS(ctx context.Context, id string)
}
