-- +migrate Up
CREATE TABLE IF NOT EXISTS gratitude_notes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    category VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);
-- +migrate Down
DROP TABLE IF EXISTS gratitude_notes; 