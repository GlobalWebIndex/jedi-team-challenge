# Jedi Team Challenge - Louk Chatwalker

---

## Description

This service provides a REST API that enables the creation of chat sessions, sending messages, reading chat sessions, 
and submitting feedback for messages. The responses to the chat messages are generated using a Retrieval-Augmented 
Generation (RAG) approach, leveraging a provided dataset of GWI data.

---

## Run

`cd script && make start-app`

* This command will start the app with `localhost` address and `:8080` port (specified in build/Dev.Dockerfile and .env)

Then you can create, get User's chat-sessions, send message and get response from the Knowledge Base (data.md)
like the examples in `/examples` directory. To generate the needed Bearer token, please call `/token` endpoint with
username = "user" & password = "password" like in the example.

This runs the app with "dlv" so that we can also attach a debugger while running.

Also you can run the `sh /scripts/e2e.sh` script to run all cases of the assignment:

1. It creates a JWT token
2. It creates three chat sessions for that User
3. In the first chat session we send three messages related to each other, so that the history
   is shown:
    1. "what do you know about latino mobile gamers?"
    2. "do they use social media?"
    3. "what social media do they use the most?"
4. It shows the whole chat session that shows the whole story
5. It submits a negative feedback for the last message
6. It sends a final message that the chat is not supposed to answer (what are butterflies)

---

## Makefile Commands

| Command                       | Usage                                            |
|-------------------------------|--------------------------------------------------|
| start-app                     | `Start app`                                      |
| kill-app                      | `Stop app`                                       |
| rebuild-app                   | `Rebuild app in case of code changes`            |
| tests                         | `Run both unit and integration tests`            |
| generate-mock FILE={filePath} | `Generate mock for a specific file`              |
| swag                          | `Generates swagger.json definitions in Docs dir` |

* All these are executed through docker containers
* In order to execute makefile commands type **make** plus a command from the table above

  make {command}

---

## Notes

1. `/config/.env` and `deployment/.env` are not pushed to Git, so in order to run the app you need them with secret keys (e.g Pinecone, OpenAI etc)
2. There are three Dockerfile files.
    1. Dockerfile is the normal, production one
    2. Dev.Dockerfile is for setting up a remote debugger Delve
    3. Utilities.Dockerfile is for building a docker for "utilities" like running tests, linting etc
4. LLM Choices (Made with limited knowledge):
    1. Pinecone for vector database
    2. `text-embedding-3-small` as embedding model
    3. Tiktoken as a tokenizer with CHUNK_ENCODING_MODEL `cl100k_base` and MAX_TOKENS_PER_CHUNKS `3000`
    4. The top 7 results are retrieved from the similarity search in the Vector DB, and there is a threshold of 0.35
       that rejects the matches with score less than that. If no such matches are found, then the answer is "The force
       is not strong enough for me to answer that question based on my context."
    5. For OpenAI model I have chosen `gpt-4.1-nano` which is a nice combination and balance of speed, accuracy and price.
    6. I've put a rate limiter when calling OpenAI because at times I was having 429 Many Request response. 
5. There are swagger definitions in `/docs`, and examples in `/examples` that show the usage of the API. And the `e2e.sh` that
   checks everything.
6. My approach for the code structure is the Hexagonal Architecture, more on that https://medium.com/@matiasvarela/hexagonal-architecture-in-go-cfd4e436faa3

## Known Issues

1. Only happy path tests are created due to time constraints.
2. JWT mechanism just requires a fake username and password to generate a JWT token and does NOT do
   actual login due to lack of time. Also no test created for it. Also, the user_id that exists in the endpoint should
   come the JWT directly.
3. In my implementation, when inputing the Chat History from the Messages DB, I import all messages to OpenAI so that
   the discussion gets continued. In a production env, I would not do that, but put a limit to the number of messages read
   from history, as there might be a lot of messages.
4. For performance increase, we can put indices in the DB, on the foreign keys so that the fetch in the GET
   endpoints is faster.

## Security

1. JWT mechanism added for Authentication and Authorization (incomplete - see Known Issues)

## Libraries and Tools

1. github.com/gorilla/mux for routing
2. gorm.io/gorm as ORM for my PostgreSQL DB
3. github.com/golang-jwt/jwt/v4 for the JWT token handling
4. github.com/openai/openai-go for communicating with OpenAI
5. github.com/pinecone-io/go-pinecone/v3 for Pinecone Vector DB
6. github.com/pkoukk/tiktoken-go for tokenizer
7. github.com/swaggo/swag for Swagger definition
8. github.com/stretchr/testify & go.uber.org/mock & github.com/DATA-DOG/go-sqlmock for testing

## Example Chat Session Response

```
{
  "id": "5488a398-1801-4a7c-ba6d-69d833453313",
  "title": "Latino Mobile Gamers Overview",
  "createdAt": "2025-05-29 10:13:58.368579 +0000 UTC",
  "updatedAt": "2025-05-29 10:13:59.831385 +0000 UTC",
  "messages": [
    {
      "id": "7a93a664-17fa-40a4-bd59-7d2c81d7accf",
      "sender": "USER",
      "content": "what do you know about latino mobile gamers",
      "created_at": "2025-05-29 10:13:58.417981 +0000 UTC"
    },
    {
      "id": "3bf684be-39ad-4fd0-8cae-27d42ce56111",
      "sender": "SYSTEM",
      "content": "Based on the provided context, Latino mobile gamers are characterized by the following behaviors and interests:\n- They are 49% more likely to visit Reddit daily compared to the average person.\n- They are 51% more likely to use TikTok weekly compared to the average person.\n- They are 42% more likely to use TikTok daily compared to the average person.\n- They are 36% more likely to use Instagram more than once a day than the average person.\n- They are 74% more likely to be interested in Esports compared to the average person.\n- They are 49% more likely to buy products or services to access the community built around them.\n- They are 62% more likely to find out about new brands and products through vlogs.\n- They are 43% more likely to discover new brands through ads seen in video or mobile games.\n- They are 62% more likely to discover new brands and products through posts or reviews from expert bloggers.\n- They are 25% between the ages of 16 and 24, and 22% between 25 and 34.\n- They are 30% more likely to be between 25 and 34 years old than the average person.\n- They are 16% of U.S. Hispanic/Latino mobile gamers.\n- They are 42% more likely to be interested in computers and coding.\n- They are 74% more likely to discover new brands and products through vlogs.\n- They are 103% more likely to discover new brands through ads in video or mobile games compared to the average person.\n- They are 45% more likely to spend 2-3 hours on streaming services daily.\n- They are 41% more likely to spend more than 4 hours on streaming services daily.\nThis indicates that Latino mobile gamers are highly active on social media platforms like TikTok and Instagram, have a strong interest in gaming, esports, technology, and discovering new brands through video content and ads in mobile games.",
      "created_at": "2025-05-29 10:14:11.136485 +0000 UTC"
    },
    {
      "id": "86718efc-d4b7-46b6-ac3f-b8c1a044cd8d",
      "sender": "USER",
      "content": "do they use social media",
      "created_at": "2025-05-29 10:14:11.172935 +0000 UTC"
    },
    {
      "id": "4dbba2a9-8df5-4c29-97f9-dd7d8999a34e",
      "sender": "SYSTEM",
      "content": "Yes, based on the provided context, Latino mobile gamers actively use social media. They are more likely than the average person to use platforms such as TikTok (42% more likely to use weekly, 36% more likely to use daily), Instagram (62% more likely to use more than once a day), and Reddit (49% more likely to visit daily). They also frequently discover new brands and products through social media content like vlogs and reviews from bloggers, indicating high engagement with social media channels.",
      "created_at": "2025-05-29 10:14:17.505033 +0000 UTC"
    },
    {
      "id": "e63bbd20-0519-45bd-a88f-48e0d4179c9e",
      "sender": "USER",
      "content": "what social media do they use the most",
      "created_at": "2025-05-29 10:14:17.541732 +0000 UTC"
    },
    {
      "id": "28e0d999-19b5-40d7-a23c-b31012fa8c87",
      "sender": "SYSTEM",
      "content": "Based on the provided context, Latino mobile gamers use TikTok and Instagram the most. They are 42% more likely to use TikTok weekly, 36% more likely to use it more than once a day, and 62% more likely to use Instagram more than once a day compared to the average person.",
      "created_at": "2025-05-29 10:14:23.200825 +0000 UTC"
    }
  ]
}```