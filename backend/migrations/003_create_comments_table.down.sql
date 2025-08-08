-- Drop indexes
DROP INDEX IF EXISTS idx_comments_created_at;
DROP INDEX IF EXISTS idx_comments_author_id;
DROP INDEX IF EXISTS idx_comments_post_id;

-- Drop comments table
DROP TABLE IF EXISTS comments;