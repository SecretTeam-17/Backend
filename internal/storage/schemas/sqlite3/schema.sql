--
-- SQL транзакция создает пустые таблицы.
--

CREATE TABLE IF NOT EXISTS modules (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS questions (
    id INTEGER PRIMARY KEY,
    module INTEGER NOT NULL REFERENCES modules(id),
    text TEXT NOT NULL,
    has_input BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS answers (
    id INTEGER PRIMARY KEY,
    question_id INTEGER NOT NULL REFERENCES questions(id),
    text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE CHECK(email LIKE '%@%')
);

CREATE INDEX IF NOT EXISTS idx_email ON users(email);

CREATE TABLE IF NOT EXISTS game_sessions (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    current_module INTEGER REFERENCES modules(id) DEFAULT 1,
    completed BOOLEAN DEFAULT FALSE,
    any_field_one TEXT DEFAULT '',
    any_field_two TEXT DEFAULT '',
    modules TEXT NOT NULL DEFAULT '',
    minigame TEXT NOT NULL DEFAULT ''
);
