DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'chat_type') THEN
        CREATE TYPE chat_type AS ENUM ('GLOBAL', 'DIRECT', 'GROUP');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS chats (
    chat_id UUID PRIMARY KEY,
    type chat_type NOT NULL,
    title TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);