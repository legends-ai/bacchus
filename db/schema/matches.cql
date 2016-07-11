CREATE TABLE IF NOT EXISTS athena.matches (
  id text,
  region text,
  body text,
  rank bigint,
  patch text,
  PRIMARY KEY(id, rank)
);
