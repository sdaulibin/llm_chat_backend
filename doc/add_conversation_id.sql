-- 向messages表添加conversation_id字段
ALTER TABLE messages ADD COLUMN conversation_id uuid NULL;