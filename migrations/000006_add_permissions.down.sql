-- ALTER TABLE users_permissions DROP CONSTRAINT IF EXISTS user_permissions_permission_id_fkey;
-- ALTER TABLE permissions DROP CONSTRAINT IF EXISTS permissions_permission_id_fkey;
-- ALTER TABLE users DROP CONSTRAINT IF EXISTS users_permission_id_fkey;
DROP TABLE IF EXISTS users_permissions cascade;;
DROP TABLE IF EXISTS permissions cascade;;
