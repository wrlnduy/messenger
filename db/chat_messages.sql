CREATE TABLE IF NOT EXISTS chat_messages (
    message_id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    text TEXT NOT NULL,
    timestamp BIGINT NOT NULL
);
