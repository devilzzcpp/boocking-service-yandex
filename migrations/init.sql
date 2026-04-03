CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    place_id INTEGER NOT NULL,
    time_from TIMESTAMP NOT NULL,
    time_to TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_bookings_place_time
    ON bookings (place_id, time_from, time_to);

CREATE INDEX IF NOT EXISTS idx_bookings_user_time
    ON bookings (user_id, time_from, id);
