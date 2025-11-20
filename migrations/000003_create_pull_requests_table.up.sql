-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE pull_requests (
    pull_request_id VARCHAR(255) PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'OPEN',
    reviewer1_id VARCHAR(255),
    reviewer2_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    merged_at TIMESTAMPTZ
);

CREATE INDEX idx_pull_requests_status ON pull_requests(status);
CREATE INDEX idx_pull_requests_author ON pull_requests(author_id);
CREATE INDEX idx_pull_requests_reviewer1 ON pull_requests(reviewer1_id);
CREATE INDEX idx_pull_requests_reviewer2 ON pull_requests(reviewer2_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE pull_requests;