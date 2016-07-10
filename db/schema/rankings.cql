CREATE TYPE IF NOT EXISTS athena.ranking (
  time timestamp,
  rank bigint
);

CREATE TABLE IF NOT EXISTS athena.rankings (
  id text,
  set<athena.ranking> rankings,
  PRIMARY KEY(id)
);
