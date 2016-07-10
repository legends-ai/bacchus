CREATE TABLE athena.matches (
  id text,
  match_id int,
  region text,
  body text,
  rank bigint,
  patch text,
  PRIMARY KEY((id, patch), rank)
);
