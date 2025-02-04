package interfaces

type PgRepositoryInterface interface {
	Ping() error
}
