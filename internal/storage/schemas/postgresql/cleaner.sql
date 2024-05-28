--
-- SQL транзакция удаляет таблицы ниже и создает их заново. Затем наполняет некоторыми тестовыми данными.
-- Работает в постгресе. 
--

-- !!! Не соответствует текущим таблицам.

DROP TABLE IF EXISTS users, game_sessions, answers, questions, modules;

CREATE TABLE modules (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT DEFAULT ''
);

CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    module INTEGER REFERENCES modules(id) DEFAULT 0,
    has_input BOOLEAN DEFAULT FALSE,
    text TEXT NOT NULL
);

CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    question_id INTEGER REFERENCES questions(id) DEFAULT 0,
    text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE CHECK(email LIKE '%@%')
);

CREATE TABLE game_sessions (
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

TRUNCATE TABLE game_sessions, users, answers, questions, modules;

INSERT INTO modules (name, description) VALUES
    ('Модуль 1', 'Описание первого модуля'),
    ('Модуль 2', 'Описание второго модуля'),
    ('Псих-модуль', 'Описание псих-модуля');

INSERT INTO modules (id, name, description) VALUES
    (0, 'Модуль 0', 'Служебный модуль');

INSERT INTO modules (name, description) VALUES
    ('Модуль 3', 'Описание третьего модуля'),
    ('Модуль 4', 'Описание четвертого модуля');

INSERT INTO questions (module, has_input, text) VALUES
    (1, false, 'Что нужно узнать у владельца перед передержкой?'),
    (1, false, 'Что нужно сделать до передержки?'),
    (2, false, 'Как правильно знакомить собак на передержке?'),
    (2, true, 'Вопрос с инпутом');

INSERT INTO questions (id, module, text) VALUES
    (0, 0, 'Служебный вопрос');

INSERT INTO answers (question_id, text, is_correct) VALUES
    (1, 'была ли ранее собака на передержке?', true),
    (1, 'есть у собаки страхи, тревожность или агрессия?', true),
    (1, 'грызет ли собака вещи дома, разрушает ли имущество когда остается одна?', true),
    (1, 'есть ли какие-то проблемы со здоровьем и аллергии?', true),
    (1, 'привита ли собака в соответствии с возрастом?', true),

    (2, 'познакомиться с владельцем и спросить подробно про собаку', true),
    (2, 'подготовить квартиру, убрать все ценные вещи повыше, а обувь в шкаф, спрятать все провода', true),
    (2, 'купить побольше еды, позвать гостей и устроить вечеринку', false),
    (2, 'прибраться перед приездом собаки', true),

    (3, 'выйти на улицу и познакомить собак на нейтральной территории, все собаки на поводке', true),
    (3, 'в квартире, все собаки без поводков', false),
    (3, 'на улице, собаки могут быть без поводков', false),

    (4, 'ответ для проверки инпута игрока', true);
