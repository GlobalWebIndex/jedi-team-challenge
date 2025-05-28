package chunks

import (
	"github.com/pkoukk/tiktoken-go"
	"strings"
)

type Chunker struct {
	Encoder           *tiktoken.Tiktoken
	MaxTokensPerChunk int
}

func NewChunker(encoder *tiktoken.Tiktoken, maxTokens int) (*Chunker, error) {
	return &Chunker{Encoder: encoder, MaxTokensPerChunk: maxTokens}, nil
}

func (c *Chunker) Chunk(text string) []string {
	var chunks []string
	sentences := splitToSentences(text)

	var buf strings.Builder
	var count int

	flush := func() {
		s := strings.TrimSpace(buf.String())
		if s != "" {
			chunks = append(chunks, s)
		}
		buf.Reset()
		count = 0
	}

	for _, sent := range sentences {
		tokCount := len(c.Encoder.Encode(sent, nil, nil))
		if tokCount > c.MaxTokensPerChunk {

			for _, w := range strings.Fields(sent) {
				wTok := len(c.Encoder.Encode(w+" ", nil, nil))
				if count+wTok > c.MaxTokensPerChunk {
					flush()
				}
				buf.WriteString(w + " ")
				count += wTok
			}
			continue
		}
		if count+tokCount > c.MaxTokensPerChunk {
			flush()
		}
		buf.WriteString(sent)
		count += tokCount
	}
	flush()
	return chunks
}

func splitToSentences(text string) []string {
	var sents []string
	var sb strings.Builder
	for _, ch := range text {
		sb.WriteRune(ch)
		if ch == '.' || ch == '!' || ch == '?' {
			sents = append(sents, sb.String())
			sb.Reset()
		}
	}

	if rem := strings.TrimSpace(sb.String()); rem != "" {
		sents = append(sents, rem)
	}
	return sents
}
