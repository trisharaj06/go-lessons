-- +goose Up

ALTER TABLE users
ADD COLUMN bio TEXT;

-- +goose Down

ALTER TABLE users
DROP COLUMN IF EXISTS bio;