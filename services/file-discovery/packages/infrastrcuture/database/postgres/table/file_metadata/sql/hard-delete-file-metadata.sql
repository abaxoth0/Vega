DELETE FROM file_metadata WHERE id = $1 AND deleted_at IS NOT NULL;
