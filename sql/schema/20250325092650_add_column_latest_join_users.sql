-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN latest_join TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
UPDATE users SET latest_join = CURRENT_TIMESTAMP WHERE latest_join IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN latest_join;
-- +goose StatementEnd
