
CREATE TABLE jobs (
  id   BIGSERIAL PRIMARY KEY,
  queue integer,
  payload text,
  status  text
);
