#!/bin/bash

# Define the backend API URL
BASE_URL="http://localhost:8080"
TOKEN_ENDPOINT="/token"
USER_ID="12345678-0000-0000-0000-000000000000"

# Step 1: Get the token
response_token=$(curl -s --location "$BASE_URL$TOKEN_ENDPOINT" \
  --header 'Content-Type: application/json' \
  --data '{"username": "user", "password": "password"}')

# Extract the token from the response using jq (ensure jq is installed)
token=$(echo "$response_token" | jq -r '.token')

# Check if the token was successfully retrieved
if [ -z "$token" ]; then
  echo "Error: Failed to retrieve token."
  exit 1
fi

echo "Token retrieved successfully: $token"
echo "----------------------------------"

# Step 2: Create 3 chat sessions for the user and get the first session's ID
for i in {1..3}; do
  # Make a POST request to create a chat session
  response=$(curl -s --location --request POST "$BASE_URL/users/$USER_ID/chat-sessions" \
    --header "Authorization: Bearer $token")

  # Extract the chat session ID from the response of the first chat session
  if [ $i -eq 1 ]; then
    CHAT_SESSION_ID=$(echo "$response" | jq -r '.id')
    echo "First chat session created with ID: $CHAT_SESSION_ID"
  fi

  # Check if the session creation was successful
#  if [[ $(echo "$response" | jq -r '.status') == "success" ]]; then
#    echo "Chat session $i created successfully."
#  else
#    echo "Error creating chat session $i."
#    exit
#  fi
done

# Step 3: Send the first message
message_content_1="what do you know about latino mobile gamers"
response_message_1=$(curl -s --location "$BASE_URL/users/$USER_ID/chat-sessions/$CHAT_SESSION_ID/messages" \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer $token" \
  --data "{\"content\":\"$message_content_1\"}")

# Step 4: Send the second message
message_content_2="do they use social media"
response_message_2=$(curl -s --location "$BASE_URL/users/$USER_ID/chat-sessions/$CHAT_SESSION_ID/messages" \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer $token" \
  --data "{\"content\":\"$message_content_2\"}")

# Step 4.5: Send the third message
message_content_3="what social media do they use the most"
response_message_3=$(curl -s --location "$BASE_URL/users/$USER_ID/chat-sessions/$CHAT_SESSION_ID/messages" \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer $token" \
  --data "{\"content\":\"$message_content_3\"}")

## Step 5: Display the responses from sending the messages
#echo "Message 1 sent to chat session $CHAT_SESSION_ID:"
#echo "$response_message_1"
#
#echo "Message 2 sent to chat session $CHAT_SESSION_ID:"
#echo "$response_message_2"
#
#echo "Message 3 sent to chat session $CHAT_SESSION_ID:"
#echo "$response_message_3"

# Step 6: Fetch and display the whole chat session
echo "Fetching the entire chat session $CHAT_SESSION_ID..."
chat_session_response=$(curl -s --location "$BASE_URL/chat-sessions/$CHAT_SESSION_ID" \
  --header "Authorization: Bearer $token")

# Display the full chat session
echo "Full Chat Session $CHAT_SESSION_ID:"
echo "$chat_session_response"
