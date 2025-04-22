// Package model содержит основные структуры данных, используемые в приложении
// для обработки запросов и ответов, связанных с сокращением URL.
package model

// ShortenerRequest представляет собой входной запрос на сокращение URL.
//
// Используется в теле HTTP-запроса.
type ShortenerRequest struct {
	// URL — оригинальный URL, который необходимо сократить.
	URL string `json:"url"`
}

// ShortenerResponse представляет ответ на успешное сокращение URL.
//
// Используется в теле HTTP-ответа.
type ShortenerResponse struct {
	// Result — сгенерированный короткий URL.
	Result string `json:"result"`
}

// ShortenURL представляет полную запись о сокращённой ссылке,
// включая уникальный ID, короткий URL и оригинальный URL.
//
// Используется как основная структура при сериализации/десериализации.
type ShortenURL struct {
	// UUID — уникальный числовой идентификатор записи.
	UUID int `json:"uuid"`

	// ShortURL — сгенерированная короткая ссылка.
	ShortURL string `json:"short_url"`

	// OriginalURL — исходный URL, который был сокращён.
	OriginalURL string `json:"original_url"`
}

// ShortenerURLMapping используется для массовой обработки сокращений.
// Содержит информацию о корреляции (например, ID клиента) и оригинальный URL.
type ShortenerURLMapping struct {
	// CorrelationID — произвольный ID, привязанный к запросу клиента.
	CorrelationID string `json:"correlation_id"`

	// OriginalURL — оригинальный URL, подлежащий сокращению.
	OriginalURL string `json:"original_url"`
}

// ShortenerURLResponse представляет результат сокращения,
// возвращаемый для каждого элемента при массовой обработке.
type ShortenerURLResponse struct {
	// CorrelationID — ID, соответствующий исходному запросу.
	CorrelationID string `json:"correlation_id"`

	// ShortURL — сгенерированная короткая ссылка.
	ShortURL string `json:"short_url"`
}

// ShortenerURLSForUserResponse описывает одну запись для выдачи пользователю,
// включающую как короткую, так и оригинальную ссылку.
type ShortenerURLSForUserResponse struct {
	// ShortURL — короткий URL, связанный с пользователем.
	ShortURL string `json:"short_url"`

	// OriginalURL — исходный URL, связанный с пользователем.
	OriginalURL string `json:"original_url"`
}
