-- 001_outreach_sends.sql
-- Tracks every email sent (or attempted) from the outreach swipe deck.

CREATE TABLE IF NOT EXISTS outreach_sends (
    id              bigserial PRIMARY KEY,
    outreach_id     bigint NOT NULL REFERENCES outreach(id) ON DELETE CASCADE,
    founder_id      bigint,
    recipient_email text NOT NULL,
    subject         text,
    message         text NOT NULL,
    status          text NOT NULL DEFAULT 'sent',
    error           text,
    user_id         text,
    sent_at         timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_outreach_sends_outreach_id ON outreach_sends (outreach_id);
CREATE INDEX IF NOT EXISTS idx_outreach_sends_founder_id  ON outreach_sends (founder_id);
CREATE INDEX IF NOT EXISTS idx_outreach_sends_user_id     ON outreach_sends (user_id);
CREATE INDEX IF NOT EXISTS idx_outreach_sends_sent_at     ON outreach_sends (sent_at);
