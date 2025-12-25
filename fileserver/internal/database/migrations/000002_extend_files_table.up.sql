ALTER TABLE files ADD COLUMN thumbnail_path TEXT;
ALTER TABLE files ADD COLUMN thumbnail_status TEXT NOT NULL DEFAULT 'pending';
DROP INDEX IF EXISTS idx_file_meta_thumbnail_status;

UPDATE files 
SET thumbnail_status = 'not_applicable' 
WHERE mime NOT LIKE 'image/%';

CREATE INDEX idx_file_meta_thumbnail_status ON files(thumbnail_status);