-- 002_saved_jobs.sql
-- Stores jobs a user has saved via the "Save job" button.

CREATE TABLE IF NOT EXISTS saved_jobs (
    id         bigserial PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    job_id     bigint NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT saved_jobs_user_job_unique UNIQUE (user_id, job_id)
);

CREATE INDEX IF NOT EXISTS idx_saved_jobs_user_id ON saved_jobs (user_id);
CREATE INDEX IF NOT EXISTS idx_saved_jobs_job_id  ON saved_jobs (job_id);
