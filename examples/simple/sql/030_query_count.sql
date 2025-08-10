-- A simple SELECT to count rows
SELECT COUNT(*) AS row_count
FROM `{{ .project_id }}`.`{{ .dataset }}`.events;


