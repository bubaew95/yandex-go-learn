package interfaces

type ShortenerRepositoryInterface interface {
	GetURLById(id string) (string, bool)
	SetURL(id string, url string)
	GetBaseURL() string
}
