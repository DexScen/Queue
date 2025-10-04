CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    max_slots INT NOT NULL CHECK (max_slots > 0),
    duration_seconds INT NOT NULL DEFAULT 600 CHECK (duration_seconds > 0)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user'
        CHECK (role IN ('user', 'admin'))
);

CREATE TABLE IF NOT EXISTS queue (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    game_id INT NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    position INT NOT NULL CHECK (position >= 0),
    joined_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL DEFAULT 'waiting'
        CHECK (status IN ('waiting', 'active', 'skipped', 'finished'))
);

CREATE INDEX idx_queue_game ON queue(game_id, position);
CREATE INDEX idx_queue_user ON queue(user_id);

-- === Игры ===
INSERT INTO games (name, description, max_slots, duration_seconds)
VALUES
    ('Code Challenge', 'Мини-задачи по программированию', 3, 600),
    ('VR Racing', 'Гонки в VR-шлемах', 2, 300),
    ('Quiz Battle', 'Интерактивная викторина', 4, 420);

-- === Пользователи ===
INSERT INTO users (login, password_hash, role)
VALUES
    ('alice', 'hash1', 'user'),
    ('bob', 'hash2', 'user'),
    ('charlie', 'hash3', 'user'),
    ('diana', 'hash4', 'user'),
    ('test', '$2a$10$DK1jX0h4oMMfezmSyf43FeEnabdqBO5kSVoXtFRxaE3Qa047Gctlm', 'user'),
    ('edward', 'hash5', 'user');

-- === Очередь ===
INSERT INTO queue (user_id, game_id, position, status)
VALUES (1, 1, 1, 'waiting');

INSERT INTO queue (user_id, game_id, position, status)
VALUES (1, 2, 1, 'waiting');

INSERT INTO queue (user_id, game_id, position, status)
VALUES (5, 2, 2, 'waiting');

INSERT INTO queue (user_id, game_id, position, status)
VALUES (5, 3, 1, 'waiting');

INSERT INTO queue (user_id, game_id, position, status)
VALUES (5, 1, 2, 'waiting');

INSERT INTO queue (user_id, game_id, position, status)
VALUES (2, 1, 3, 'waiting');