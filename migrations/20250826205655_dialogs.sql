-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dialogs (
    id SERIAL PRIMARY KEY,
    from_user_id VARCHAR(255) NOT NULL,
    to_user_id VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dialogs;
-- +goose StatementEnd
