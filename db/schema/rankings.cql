CREATE TYPE IF NOT EXISTS athena.ranking (
  time timestamp,
  rank bigint
);

CREATE TABLE IF NOT EXISTS athena.rankings (
  id text,
  rankings set<frozen<athena.ranking>>,
  rank bigint,
  PRIMARY KEY(id, rank)
);
