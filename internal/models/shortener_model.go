package models

type ShortenerRequest struct {
	URL string `json:"url"`
}

type ShortenerResponse struct {
	Result string `json:"result"`
}

type ShortenURL struct {
	UUID        int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
