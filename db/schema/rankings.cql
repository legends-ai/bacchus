-- ranking is a rank and its fetch time.
CREATE TYPE IF NOT EXISTS athena.ranking (
  time timestamp,
  rank bigint
)

-- rankings stores all rankings.
CREATE TABLE athena.rankings (
  -- id is the summoner id.
  id text,

  -- ranks is a series of rankings throughout time.
  set<athena.ranking> rankings,

  PRIMARY KEY(id)
)
