package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/loukaspe/jedi-team-challenge/pkg/chunks"
	"github.com/loukaspe/jedi-team-challenge/pkg/embeddings"
	"github.com/loukaspe/jedi-team-challenge/pkg/vectordb"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/pinecone-io/go-pinecone/v3/pinecone"
	"github.com/pkoukk/tiktoken-go"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func main() {
	getEnv()

	pineconeAPIKey := os.Getenv("PINECONE_API_KEY")
	pineconeIndexName := os.Getenv("PINECONE_INDEX")
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	//openAIURL := os.Getenv("OPENAI_URL")
	chunkEncoding := os.Getenv("CHUNK_ENCODING_MODEL")
	embeddingModel := openai.EmbeddingModel(os.Getenv("EMBEDDING_MODEL"))
	maxTokensPerChunksAsString := os.Getenv("MAX_TOKENS_PER_CHUNKS")
	maxTokensPerChunks, err := strconv.Atoi(maxTokensPerChunksAsString)
	if err != nil {
		log.Fatal("Cannot read max token per chunks: ", err)
	}

	tiktokenEncoder, err := tiktoken.GetEncoding(chunkEncoding)
	if err != nil {
		log.Fatal("Cannot create encoder: ", err)
	}

	chunker, err := chunks.NewTiktokenChunker(tiktokenEncoder, maxTokensPerChunks)
	if err != nil {
		log.Fatal("Cannot create chunker: ", err)
	}

	openAIClient := openai.NewClient(option.WithAPIKey(openAIAPIKey))
	embedder := embeddings.NewEmbeddingService(&openAIClient, embeddingModel)

	textBytes, err := os.ReadFile("./data.md")
	if err != nil {
		log.Fatal(err)
	}
	text := string(textBytes)

	chunks := chunker.Chunk(text)
	fmt.Printf("Generated %d chunks\n", len(chunks))

	ctx := context.Background()
	embeddings, err := embedder.Embed(ctx, chunks)
	if err != nil {
		log.Fatalf("Embedding error: %v", err)
	}

	// Output
	for i, emb := range embeddings {
		fmt.Printf("Chunk %d → %d-dim vector\n", i, len(emb))
	}

	pineconeClient, err := pinecone.NewClient(pinecone.NewClientParams{
		ApiKey: pineconeAPIKey,
	})
	if err != nil {
		log.Fatalf("Failed to create pinecone Client: %v", err)
	}
	pineconeVectorDB := vectordb.NewPineconeVectorDB(
		pineconeIndexName,
		pineconeClient,
	)

	count, err := pineconeVectorDB.StoreEmbeddings(ctx, embeddings)
	if err != nil {
		log.Fatalf("Failed to store embeddings: %v", err)
	}

	fmt.Sprintf("Stored %d embeddings in Pinecone index %s\n", count, pineconeIndexName)

	//logger := logger.NewLogger(context.Background())
	//router := mux.NewRouter()
	//httpServer := &http.Server{
	//	Addr:    os.Getenv("SERVER_ADDR"),
	//	Handler: router,
	//}
	//db := getDB()
	//
	//server := server.NewServer(db, router, httpServer, logger)
	//
	//server.Run()
}

//func getDB() *gorm.DB {
//	dbDsn := fmt.Sprintf(
//		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s TimeZone=Europe/Athens",
//		os.Getenv("DB_HOST"),
//		os.Getenv("DB_PORT"),
//		os.Getenv("DB_USER"),
//		os.Getenv("DB_NAME"),
//		os.Getenv("DB_PASSWORD"),
//	)
//	db, err := gorm.Open(postgres.Open(dbDsn), &gorm.Config{})
//	if err != nil {
//		log.Fatal("Cannot connect to database: ", err)
//	}
//
//	// Drops added in order to start with clean DB on App start for
//	// assessment reasons
//	db.Migrator().DropTable("users")
//
//	err = db.AutoMigrate(&repositories.User{})
//	if err != nil {
//		log.Fatal("cannot migrate user table")
//	}
//
//	return db
//}

func getEnv() {
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	}
}
