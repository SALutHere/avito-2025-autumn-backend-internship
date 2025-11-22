CREATE TABLE IF NOT EXISTS teams (
    name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
    id        TEXT PRIMARY KEY,
    username  TEXT NOT NULL,
    team_name TEXT NOT NULL REFERENCES teams(name) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name);

CREATE TABLE IF NOT EXISTS pull_requests (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    author_id  TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status     TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    merged_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_pr_author ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_pr_status ON pull_requests(status);

CREATE TABLE IF NOT EXISTS pull_request_reviewers (
    pr_id       TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id TEXT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    PRIMARY KEY (pr_id, reviewer_id)
);

CREATE INDEX IF NOT EXISTS idx_reviewer_lookup ON pull_request_reviewers(reviewer_id);
