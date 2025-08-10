-- Inserts a couple of example rows, using templated vars
INSERT INTO `{{ .project_id }}`.`{{ .dataset }}`.events (event_id, event_ts, env, user)
VALUES
  ("evt-1", CURRENT_TIMESTAMP(), "{{ .env }}", SESSION_USER()),
  ("evt-2", CURRENT_TIMESTAMP(), "{{ .env }}", SESSION_USER());
