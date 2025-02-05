-- +goose Up
-- +goose StatementBegin
CREATE TABLE chirps
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP    NOT NULL,
    updated_at TIMESTAMP    NOT NULL,
    body       VARCHAR(255) NOT NULL,
    user_id    UUID         NOT NULL,
    CONSTRAINT FK_users FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chirps
DROP CONSTRAINT FK_users; 
DROP TABLE chirps;
-- +goose StatementEnd
