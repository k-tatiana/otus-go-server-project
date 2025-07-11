-- +goose Up
-- +goose StatementBegin
CREATE TYPE gender AS ENUM ('male', 'female', 'undefined');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE gender;
-- +goose StatementEnd
