CREATE TABLE IF NOT EXISTS chat_messages (
    message_id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id),
    
    text TEXT NOT NULL,
    timestamp BIGINT NOT NULL
);
