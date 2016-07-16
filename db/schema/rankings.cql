CREATE TABLE IF NOT EXISTS athena.rankings (
  id text,
  time timestamp,
  rank bigint,
  PRIMARY KEY(id, time, rank)
);
