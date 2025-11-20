-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    team_name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviewers (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id),
    github_username VARCHAR(255) NOT NULL UNIQUE,
    capacity INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE reviewers;
DROP TABLE teams;