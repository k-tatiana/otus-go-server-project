-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    birthday DATE NOT NULL,
    gender gender,
    interests TEXT[],
    city VARCHAR(255),
    login VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    token uuid DEFAULT gen_random_uuid() PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
