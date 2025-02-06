package interfaces

type ShortenerServiceInterface interface {
	GenerateURL(url string, randomStringLength int) string
	RandStringBytes(n int) string
	GetURLByID(id string) (string, bool)
	GetAllURL() map[string]string
	Ping() error
}
