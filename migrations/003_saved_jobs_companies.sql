-- 003_saved_jobs_companies.sql
-- Extends saved_jobs so a user can save either a job or a company.
-- When job_id is NULL and company_id is set, the row represents a saved company.

-- job_id is now optional; a row either points to a saved job or a saved company.
ALTER TABLE saved_jobs ALTER COLUMN job_id DROP NOT NULL;

-- company_id is required on every row (saved job rows store the job's company).
ALTER TABLE saved_jobs ADD COLUMN company_id bigint REFERENCES companies(id) ON DELETE CASCADE;

-- Backfill company_id for existing saved-job rows.
UPDATE saved_jobs sj
SET company_id = j.company_id
FROM jobs j
WHERE j.id = sj.job_id;

ALTER TABLE saved_jobs ALTER COLUMN company_id SET NOT NULL;

-- The old unique constraint no longer works because job_id is now nullable
-- (Postgres treats NULLs as distinct, so (user_id, NULL) rows would never conflict).
ALTER TABLE saved_jobs DROP CONSTRAINT saved_jobs_user_job_unique;

-- A user can save a given job only once.
CREATE UNIQUE INDEX saved_jobs_user_job_unique
    ON saved_jobs (user_id, job_id)
    WHERE job_id IS NOT NULL;

-- A user can save a given company only once.
CREATE UNIQUE INDEX saved_jobs_user_company_unique
    ON saved_jobs (user_id, company_id)
    WHERE job_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_saved_jobs_company_id ON saved_jobs (company_id);
