# Example: simple

This is a minimal BigQuery SQL pipeline runnable with `bqdo`.

## Prerequisites

- Google Cloud auth set up (e.g., `gcloud auth application-default login`) or set `GOOGLE_APPLICATION_CREDENTIALS`.
- BigQuery API enabled in your project.

## Files

- `bqdo.toml`: Pipeline configuration and template variables.
- `sql/001_create_dataset.sql`: Creates the dataset.
- `sql/010_create_table.sql`: Creates a demo table.
- `sql/020_insert_sample.sql`: Inserts sample rows.
- `sql/030_query_count.sql`: Queries the count of rows.

## Run

From the repo root:

```sh
mise run build
bin/bqdo run -c examples/simple/bqdo.toml -p YOUR_PROJECT_ID -d bqdo_example -l US
```

To validate without executing:

```sh
bin/bqdo run -c examples/simple/bqdo.toml -p YOUR_PROJECT_ID -d bqdo_example -l US --dry-run
```
