# Services

| Service         | Responsibilities                                                 | Owns                                                                            |
| --------------- | ---------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| File Repository | All file CRUD operations + serving file content. Source of truth | Files permissions, quotas, data consistency                                     |
| Auth            | Authentication and authorization                                 | Users and services permissions. Essential users data (logins, passwords, roles) |
