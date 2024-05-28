--
-- SQL транзакция создает таблицы и некоторые служебные строки.
--

-- !!! Не соответствует текущим таблицам.

CREATE TABLE IF NOT EXISTS modules (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    module INTEGER REFERENCES modules(id) DEFAULT 0,
    has_input BOOLEAN DEFAULT FALSE,
    text TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS answers (
    id SERIAL PRIMARY KEY,
    question_id INTEGER REFERENCES questions(id) DEFAULT 0,
    text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE CHECK(email LIKE '%@%')
);

CREATE TABLE IF NOT EXISTS game_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    current_module INTEGER REFERENCES modules(id) DEFAULT 1,
    completed BOOLEAN DEFAULT FALSE,
    any_field_one TEXT DEFAULT '',
    any_field_two TEXT DEFAULT '',
    modules JSONB NOT NULL DEFAULT '{}'::jsonb,
    minigame JSONB NOT NULL DEFAULT '{}'::jsonb
);

INSERT INTO modules (id, name, description) VALUES
    (0, 'Модуль 0', 'Служебный модуль');

INSERT INTO questions (id, module, text) VALUES
    (0, 0, 'Служебный вопрос');

INSERT INTO users (id, name, email) VALUES
    (0, 'Служебный юзер', 'Служебный емэйл');