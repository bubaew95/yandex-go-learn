package interfaces

type ShortenerServiceInterface interface {
	GenerateID(url string, randomStringLength int) string
	RandStringBytes(n int) string
	GenerateResponseURL(id string) string
	GetURLByID(id string) (string, bool)
}
