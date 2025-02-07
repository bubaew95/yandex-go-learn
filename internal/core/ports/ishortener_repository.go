package ports

type ShortenerRepositoryInterface interface {
	GetURLByID(id string) (string, bool)
	SetURL(id string, url string)
	GetAllURL() map[string]string
	Ping() error
	Close() error
}
