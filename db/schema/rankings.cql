CREATE TABLE IF NOT EXISTS athena.rankings (
  id text,
  rank bigint,
  time timestamp,
  PRIMARY KEY(id, rank)
);
