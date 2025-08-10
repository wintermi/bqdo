-- Creates a demo table in the dataset
CREATE TABLE IF NOT EXISTS `{{ .project_id }}`.`{{ .dataset }}`.events (
  event_id STRING,
  event_ts TIMESTAMP,
  env STRING
);


