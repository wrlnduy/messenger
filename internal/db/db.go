package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		user_id UUID PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,

		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		is_admin BOOLEAN NOT NULL DEFAULT FALSE,

		created_at TIMESTAMP NOT NULL DEFAULT now()
	);

	CREATE TABLE IF NOT EXISTS chat_messages (
		message_id UUID PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(user_id),
		
		text TEXT NOT NULL,
		timestamp TIMESTAMP NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sessions (
		session_id UUID PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(user_id),
		
		expires_at TIMESTAMP NOT NULL
	);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
