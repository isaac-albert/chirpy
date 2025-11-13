-- +goose Up
CREATE TABLE userchirps (
    Id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    body TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users (Id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE userchirps;