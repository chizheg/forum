CREATE TABLE chat_messages (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_chat_messages_created_at ON chat_messages(created_at);

CREATE TABLE chat_participants (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
); 