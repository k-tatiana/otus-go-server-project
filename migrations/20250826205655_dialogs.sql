-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS dialogs (
    from_user_id VARCHAR(255) NOT NULL,
    to_user_id VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    user_pair_hash int primary key,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
SELECT create_distributed_table('dialogs', 'user_pair_hash');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dialogs;
-- +goose StatementEnd
