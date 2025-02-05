-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_token
(
    token      TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id    UUID NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT FK_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE refresh_tokens
DROP CONSTRAINT FK_users; 
DROP TABLE refresh_tokens
-- +goose StatementEnd
