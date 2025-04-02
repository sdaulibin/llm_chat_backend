-- 向messages表添加conversation_id字段
ALTER TABLE messages ADD COLUMN conversation_id uuid NULL;
-- 向messages表添加file_id字段
ALTER TABLE messages ADD COLUMN file_id uuid NULL;
-- 向messages表添加current_id和is_stop字段
ALTER TABLE messages ADD COLUMN current_id uuid NULL;
ALTER TABLE messages ADD COLUMN is_stop boolean DEFAULT false NOT NULL;
ALTER TABLE messages ADD COLUMN task_id uuid NULL;