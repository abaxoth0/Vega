BEGIN;
    CREATE TABLE IF NOT EXISTS "file_metadata" (
        id              uuid            PRIMARY KEY,
        path            VARCHAR(1024)   NOT NULL,
        bucket          uuid            NOT NULL,
        size            BIGINT          NOT NULL,
        mime            VARCHAR(128)    NOT NULL,
        permissions     SMALLINT        NOT NULL,
        status          CHAR(1)         NOT NULL DEFAULT 'A', -- A - Active; R - Archived; P - Pending
        owner           uuid            NOT NULL,
        original_name   VARCHAR(255)    NOT NULL,
        encoding        VARCHAR(8)      NOT NULL,
        categories      VARCHAR(32)[],
        tags            VARCHAR(32)[],
        checksum        CHAR(64)        NOT NULL, -- SHA256
        uploaded_by     uuid            NOT NULL,
        created_at      TIMESTAMP       NOT NULL DEFAULT NOW(),
        uploaded_at     TIMESTAMP       NOT NULL DEFAULT NOW(),
        updated_at      TIMESTAMP,
        accessed_at     TIMESTAMP,
        deleted_at      TIMESTAMP,
        description     TEXT,
    );
COMMIT;
