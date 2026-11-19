CREATE TABLE IF NOT EXISTS approvals (
    approval_id TEXT PRIMARY KEY,
    task_id TEXT,
    session_id TEXT,
    tool_name TEXT NOT NULL,
    scope_hash TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TEXT NOT NULL,
    decided_at TEXT,
    consumed_at TEXT
);

CREATE INDEX IF NOT EXISTS idx_approvals_task_id
ON approvals(task_id);

CREATE INDEX IF NOT EXISTS idx_approvals_status
ON approvals(status);

CREATE INDEX IF NOT EXISTS idx_approvals_created_at
ON approvals(created_at);
