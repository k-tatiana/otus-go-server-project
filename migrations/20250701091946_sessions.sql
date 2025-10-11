-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    token VARCHAR NOT NULL,
    expiration_time TIMESTAMP WITHOUT TIME ZONE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd
