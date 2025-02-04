package interfaces

type PgServiceInterface interface {
	Ping() error
}
