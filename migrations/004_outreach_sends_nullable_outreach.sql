-- 004_outreach_sends_nullable_outreach.sql
-- Allow sends without an outreach: cards may have no outreach, in which case
-- the send is recorded against the founder id only.

ALTER TABLE outreach_sends ALTER COLUMN outreach_id DROP NOT NULL;
