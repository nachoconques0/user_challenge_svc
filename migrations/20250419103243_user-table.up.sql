BEGIN;

CREATE TABLE challenge.user (
  id UUID PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  nickname TEXT NOT NULL,
  password TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  country TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL,
  deleted_at TIMESTAMPTZ
);

-- Update triggers
CREATE TRIGGER set_updated_at
  BEFORE INSERT OR UPDATE ON challenge.user
  FOR EACH ROW
  EXECUTE PROCEDURE challenge.set_updated_at ();

COMMIT;
