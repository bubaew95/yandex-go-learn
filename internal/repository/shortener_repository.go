package repository

type ShortenerRepository struct {
	data    map[string]string
	baseURL string
}

func NewShortenerRepository(data map[string]string, baseURL string) *ShortenerRepository {
	return &ShortenerRepository{
		data:    data,
		baseURL: baseURL,
	}
}

func (s ShortenerRepository) SetURL(id string, url string) {
	s.data[id] = url
}

func (s ShortenerRepository) GetURLByID(id string) (string, bool) {
	url, ok := s.data[id]

	return url, ok
}

func (s ShortenerRepository) GetBaseURL() string {
	return s.baseURL
}

func (s ShortenerRepository) GetAllURL() map[string]string {
	return s.data
}
