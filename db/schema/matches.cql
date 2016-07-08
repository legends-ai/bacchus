-- matches stores all match data.
CREATE TABLE athena.matches (
  -- id is the combined identifier of the match.
  id text,

  -- match_id is the Riot match id
  match_id int,

  -- region is the region (e.g NA)
  region text,

  -- body is the entire response returned from Riot.
  body text,

  -- rank is a bigint derived from division and tier.
  rank bigint,

  -- patch is the patch this game happened on.
  patch text,

  -- Cluster by patch and rank to allow sorting
  PRIMARY KEY((id, patch), rank)
)
