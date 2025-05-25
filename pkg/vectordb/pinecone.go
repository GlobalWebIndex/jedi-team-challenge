package vectordb

import (
	"context"
	"fmt"
	"github.com/pinecone-io/go-pinecone/v3/pinecone"
)

type PineconeVectorDB struct {
	client *pinecone.Client
	index  string
}

func NewPineconeVectorDB(index string, client *pinecone.Client) *PineconeVectorDB {
	return &PineconeVectorDB{
		index:  index,
		client: client,
	}
}

func (db *PineconeVectorDB) StoreEmbeddings(ctx context.Context, embeddings [][]float64) (int, error) {
	idx, err := db.client.DescribeIndex(ctx, db.index)
	if err != nil {
		return 0, err
	}

	idxConnection, err := db.client.Index(pinecone.NewIndexConnParams{Host: idx.Host})
	if err != nil {
		return 0, err
	}

	vectors := make([]*pinecone.Vector, len(embeddings))
	for i, vec := range embeddings {
		id := fmt.Sprintf("doc1-chunk-%d", i)

		//md, err := structpb.NewStruct(map[string]interface{}{
		//	"text": chunks[i],
		//})
		//if err != nil {
		//	log.Fatalf("failed to create metadata struct: %v", err)
		//}

		vectorToFloat32 := make([]float32, len(vec))
		for i, v := range vec {
			vectorToFloat32[i] = float32(v)
		}

		vectors[i] = &pinecone.Vector{
			Id:     id,
			Values: &vectorToFloat32,
			//Metadata: map[string]interface{}{
			//	"text": chunks[i], // store the original chunk if you want
			//},
		}
	}

	count, err := idxConnection.UpsertVectors(ctx, vectors)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
