-- ranking is a rank and its fetch time.
CREATE TYPE IF NOT EXISTS athena.ranking (
  time timestamp,
  rank bigint
)

-- matches stores all match data.
CREATE TABLE athena.ranks (
  -- id is the summoner id.
  id text,

  -- ranks is a series of rankings throughout time.
  set<athena.ranking> ranks,

  PRIMARY KEY(id)
)
