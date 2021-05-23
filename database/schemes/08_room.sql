CREATE TABLE IF NOT EXISTS room
 (
 room_id INTEGER NOT NULL,
 user_id INTEGER NOT NULL,
 FOREIGN KEY (user_id) REFERENCES users (id),
 FOREIGN KEY (room_id) REFERENCES rooms (id)
);