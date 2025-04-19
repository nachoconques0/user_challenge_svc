BEGIN;

CREATE SCHEMA IF NOT EXISTS challenge;

CREATE OR REPLACE FUNCTION challenge.set_updated_at ()
    RETURNS TRIGGER STABLE
    AS $plpgsql$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$plpgsql$
LANGUAGE plpgsql;

COMMIT;
