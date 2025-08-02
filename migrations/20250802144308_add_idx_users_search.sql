-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_users_search ON users USING btree (lower(name), lower(surname) varchar_pattern_ops);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index idx_users_search;
-- +goose StatementEnd
