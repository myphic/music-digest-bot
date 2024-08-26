-- +goose Up
-- +goose StatementBegin
ALTER TABLE digest
    ADD COLUMN posted BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE digest
DROP COLUMN posted;
-- +goose StatementEnd