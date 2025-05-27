package embeddings

import (
	"context"
	"errors"
	"fmt"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	"github.com/openai/openai-go"
	"math"
	"time"
)

const maxEmbeddingRetries = 5

type EmbeddingService struct {
	OpenAIClient   *openai.Client
	EmbeddingModel openai.EmbeddingModel
}

func NewEmbeddingService(client *openai.Client, embeddingModel openai.EmbeddingModel) *EmbeddingService {
	return &EmbeddingService{
		OpenAIClient:   client,
		EmbeddingModel: embeddingModel,
	}
}

func (s *EmbeddingService) Embed(ctx context.Context, inputs []string) ([]*domain.Embeddings, error) {
	var resp *openai.CreateEmbeddingResponse
	var err error
	domainEmbeddings := make([]*domain.Embeddings, 0)

	for _, input := range inputs {
		for attempt := 0; attempt < maxEmbeddingRetries; attempt++ {
			resp, err = s.OpenAIClient.Embeddings.New(
				ctx,
				openai.EmbeddingNewParams{
					Model: s.EmbeddingModel,
					Input: openai.EmbeddingNewParamsInputUnion{OfString: openai.String(input)},
				},
			)

			if err == nil {
				break
			}

			// see if it's an APIError with HTTP 429
			var apiErr *openai.Error
			if errors.As(err, &apiErr) && apiErr.StatusCode == 429 {
				wait := time.Duration(math.Pow(2, float64(attempt))) * time.Second
				// TODO: log
				fmt.Printf("Rate‐limited on attempt %d: waiting %s before retry…", attempt+1, wait)
				time.Sleep(wait)
				continue
			}
			// some other error – give up immediately
			return nil, fmt.Errorf("failed embedding: %w", err)
		}

		if err != nil {
			return nil, err
		}

		for _, d := range resp.Data {
			domainEmbeddings = append(domainEmbeddings, &domain.Embeddings{
				Embeddings: d.Embedding,
				Text:       input,
			})
		}
	}

	return domainEmbeddings, nil
}
