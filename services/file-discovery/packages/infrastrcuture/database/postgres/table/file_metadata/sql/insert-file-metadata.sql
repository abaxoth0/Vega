INSERT INTO file_metadata
    (id, path, bucket, size, mime, permissions, status, owner, original_name, encoding, categories, tags, checksum, uploaded_by, uploaded_at, created_at, description)
VALUES (
    $1, --  '3a153ed3-665a-43d1-a380-4067896fbb28',
    $2, --  '/some_directory/example.txt',
    $3, --  'ced077c7-08de-4237-9358-9778780e0592',
    $4, --  12,
    $5, --  'text/plain',
    $6, --  992,
    $7, --  'A',
    $8, --  'e2eeba2b-68b7-4125-a23e-0557c9a62066',
    $9, --  'unnamed.txt',
    $10, --  'UTF-8',
    $11, --  '{"example"}',
    $12, --  '{"testing"}',
    $13, --  '0e94ae36da6ff03992a57fddbdf4728b609d0d7fe6eb019fa9f1b9b5b540d835',
    $14, --  'e2eeba2b-68b7-4125-a23e-0557c9a62066',
    $15, --  NOW(),
    $16, --  NOW()
    $17 --  'First example file'
);
