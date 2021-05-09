CREATE TABLE
IF NOT EXISTS online_users
(
    user_id INTEGER NOT NULL,
    expires_at INTEGER NOT NULL,

FOREIGN KEY
(user_id) REFERENCES users
(id) ON
DELETE CASCADE
);