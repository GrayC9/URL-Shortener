package shortener

type URL struct {
	OriginalURL string `json:"original"`
	ShortURL    string `json:"short"`
}
