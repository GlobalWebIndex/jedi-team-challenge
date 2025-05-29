package domain

type Embeddings struct {
	Embeddings []float64 `json:"embedding"`
	Text       string    `json:"text"`
}
