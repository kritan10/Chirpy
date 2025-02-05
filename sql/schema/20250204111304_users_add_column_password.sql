-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN password TEXT NOT NULL DEFAULT 'unset';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN password;
-- +goose StatementEnd
