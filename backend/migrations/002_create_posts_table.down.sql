-- Drop trigger
DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;

-- Drop indexes
DROP INDEX IF EXISTS idx_posts_tags;
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_published;
DROP INDEX IF EXISTS idx_posts_author_id;

-- Drop posts table
DROP TABLE IF EXISTS posts;