UPDATE file_metadata SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL;
