CREATE TABLE IF NOT EXISTS chat_unread (
    chat_id UUID NOT NULL REFERENCES chats(chat_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    unread_count INT NOT NULL DEFAULT 0,
    PRIMARY KEY (chat_id, user_id)
);