-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN id BIGSERIAL;
CREATE TABLE posts (
	id bigserial primary key,
	author_user_id bigint, 
	"text" text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN id;
DROP TABLE posts;
-- +goose StatementEnd
