ALTER TABLE chat_messages
ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'active',
ADD COLUMN edited_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX idx_chat_messages_status ON chat_messages(status);
CREATE INDEX idx_chat_messages_deleted_at ON chat_messages(deleted_at); 