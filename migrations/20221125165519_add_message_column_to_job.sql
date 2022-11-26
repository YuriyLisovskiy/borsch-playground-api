-- +goose Up
-- +goose StatementBegin
ALTER TABLE jobs ADD COLUMN message TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jobs DROP COLUMN message;
-- +goose StatementEnd
