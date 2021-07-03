CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY,
    room_id INTEGER NOT NULL,
    author_id INTEGER NOT NULL,
    message TEXT NOT NULL,
    message_date INTEGER NOT NULL,
    read INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (author_id) REFERENCES users (id),
    FOREIGN KEY (room_id) REFERENCES rooms (id)
);