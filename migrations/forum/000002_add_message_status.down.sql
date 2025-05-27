DROP INDEX IF EXISTS idx_chat_messages_deleted_at;
DROP INDEX IF EXISTS idx_chat_messages_status;

ALTER TABLE chat_messages
DROP COLUMN status,
DROP COLUMN edited_at,
DROP COLUMN deleted_at; 