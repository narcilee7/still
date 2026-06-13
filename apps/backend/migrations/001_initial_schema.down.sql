DROP INDEX IF EXISTS idx_resonances_user_id;
DROP INDEX IF EXISTS idx_resonances_post_id;
DROP INDEX IF EXISTS idx_posts_user_id;
DROP INDEX IF EXISTS idx_posts_created_at;

DROP TABLE IF EXISTS resonances;
DROP TABLE IF EXISTS post_ai_analysis;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;
