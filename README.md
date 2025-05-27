# Jedi Team Challenge - Louk Chatwalker

---

## Description

Service that provides a REST API offering creating chat sessions, sending messages, reading the chat session and  
submitting feedback for a message.

---

## Run

`cd script && make start-app`

* This command will start the app with `localhost` address and `:8080` port (specified in build/Dev.Dockerfile and .env)

Then you can create, get User's chat-sessions, send message and get response from the Knowledge Base (data.md)
like the examples in /examples directory. To generate the needed Bearer token, please call /token endpoint with
username = "user"  
& password = "password" like in the example.

This runs the app with "dlv" so that we can also attach a debugger while running.

Also you can run the /scripts/e2e.sh script to run all cases of the assignment:

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

1. .env is not pushed to Git, so in order to run the app you need it with secret keys (e.g Pinecone, OpenAI etc)
2. There are three Dockerfile files.
    1. Dockerfile is the normal, production one
    2. Dev.Dockerfile is for setting up a remote debugger Delve
    3. Utilities.Dockerfile is for building a docker for "utilities" like running tests,  
       linting etc
4. LLM Choices I took: (This is a field that it's a bit unknown for me, so some decision were made with just a little
   studying
   and might not be the correct ones):
    1. Pinecone for vector database
    2. text-embedding-3-small as embedding model
    3. Tiktoken as a tokenizer with CHUNK_ENCODING_MODEL cl100k_base and MAX_TOKENS_PER_CHUNKS 3000
    4. The top 7 results are retrieved from the similarity search in the Vector DB, and there is a threshold of 0.35
       that rejects the matches with score less than that. If no such matches are found, then the answer is "The force
       is not strong enough for me to answer that question based on my context."
    5. For OpenAI model I have chosen "gpt-4.1-nano" which is a nice combination and balance of speed, accuracy and price.
5. There are swagger definitions in /docs, and examples in /examples that show the usage of the API. And the e2e.sh that
   checks everything.

## Known Issues

1. Only happy path tests are created due to time constraints.
2. JWT mechanism just requires a fake username and password to generate a JWT token and does NOT do
   actual login due to lack of time. Also no test created for it. Also, the user_id that exists in the endpoint should
   come
   the JWT directly.
3. In my implementation, when inputing the Chat History from the Messages DB, I import all messages to OpenAI so that
   the
   discussion gets continued. In a production env, I would not do that, but put a limit to the number of messages read
   from
   history, as there might be a lot of messages.

## Security

1. JWT mechanism added for Authentication and Authorization (incomplete - see Known Issues)