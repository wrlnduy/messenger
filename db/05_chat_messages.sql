CREATE TABLE IF NOT EXISTS chat_messages (
    message_id UUID PRIMARY KEY,
    chat_id UUID NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
    
    user_id UUID NOT NULL REFERENCES users(user_id),
    text TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_chat_messages_chat_id_timestamp
ON chat_messages (chat_id, timestamp);

