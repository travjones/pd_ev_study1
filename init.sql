CREATE TABLE data (
  data_id SERIAL PRIMARY KEY,
  completed timestamptz DEFAULT now(),
  data jsonb,
  ev_data jsonb
);