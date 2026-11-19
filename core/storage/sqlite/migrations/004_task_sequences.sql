CREATE TABLE IF NOT EXISTS task_sequences (
    day TEXT PRIMARY KEY,
    value INTEGER NOT NULL
);

INSERT INTO task_sequences (
    day,
    value
)
SELECT
    substr(task_id, 6, 8) AS day,
    MAX(
        CAST(
            substr(task_id, 15)
            AS INTEGER
        )
    ) AS value
FROM tasks
WHERE
    task_id LIKE 'task_%'
    AND length(task_id) >= 15
GROUP BY
    substr(task_id, 6, 8)
ON CONFLICT(day) DO UPDATE SET
    value = MAX(
        task_sequences.value,
        excluded.value
    );
